package pg

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"dishdash.ru/cmd/server/config"

	// driver for migration
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultAttempts = 20
	defaultTimeout  = time.Second
)

func init() {
	if !config.C.DB.AutoMigrate {
		log.Println("Database migrator is disabled")
		return
	}
	if err := MigrateDB(); err != nil {
		log.Fatal(err)
	}
}

func MigrateDB() error {
	dbUrl := config.C.DBUrl()
	var (
		attempts = defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	_, path, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("could not determine migration location")
	}
	pathToMigrationFiles := filepath.Dir(path) + "/../../../migrations"

	for attempts > 0 {
		m, err = migrate.New(fmt.Sprintf("file:%s", pathToMigrationFiles), dbUrl)
		if err == nil {
			break
		}

		log.Printf("Migrate: pgdb is trying to connect, attempts left: %d, error: %s", attempts, err.Error())
		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	err = m.Up()
	defer func() { _, _ = m.Close() }()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no changes")
		return nil
	}

	log.Printf("Migrate: up success")
	return nil
}
