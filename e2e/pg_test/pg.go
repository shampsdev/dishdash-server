package pg_test

//nolint: revive // test stub
import (
	"context"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"dishdash.ru/cmd/server/config"

	"github.com/testcontainers/testcontainers-go"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabase struct {
	DbInstance *pgxpool.Pool
	DbAddress  string
	container  testcontainers.Container
}

func SetupTestDatabase() *TestDatabase {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)

	config.C.DB.User = "root"
	config.C.DB.Password = "root"
	config.C.DB.Host = "localhost"
	config.C.DB.Database = "root"
	config.C.DB.AutoMigrate = false

	container, dbInstance, dbAddr, err := createContainer(ctx)
	if err != nil {
		log.Fatal("failed to setup test", err)
	}
	cancel()

	return &TestDatabase{
		container:  container,
		DbInstance: dbInstance,
		DbAddress:  dbAddr,
	}
}

func (tdb *TestDatabase) TearDown() {
	tdb.DbInstance.Close()
	// remove test container
	_ = tdb.container.Terminate(context.Background())
}

func createContainer(ctx context.Context) (testcontainers.Container, *pgxpool.Pool, string, error) {
	env := map[string]string{
		"POSTGRES_PASSWORD": config.C.DB.Password,
		"POSTGRES_USER":     config.C.DB.User,
		"POSTGRES_DB":       config.C.DB.Database,
	}

	port := "5432/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgis/postgis:16-3.4",
			Cmd:          []string{"postgres", "-c", "log_statement=all", "-c", "log_destination=stderr"},
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to start container: %w", err)
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to get container external port: %w", err)
	}
	portNum, _ := strconv.ParseInt(p.Port(), 10, 32)
	config.C.DB.Port = uint16(portNum)

	log.Info("postgres container ready and running at port: ", p.Port())

	time.Sleep(time.Second)

	dbAddr := fmt.Sprintf("localhost:%s", p.Port())

	db, err := pgxpool.NewWithConfig(ctx, config.C.PGXConfig())
	if err != nil {
		return container, db, dbAddr, fmt.Errorf("failed to establish database connection: %w", err)
	}

	return container, db, dbAddr, nil
}
