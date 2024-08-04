package pg

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate/v4"

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
	if err := MigrateUp(); err != nil {
		log.Fatal(err)
	}
}

func MigrateUp() error {
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
		log.Fatalf("Migrate up: error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate up: no changes")
		return nil
	}

	log.Printf("Migrate up: success")
	return nil
}

func MigrateDown() error {
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

		log.Printf("Migrate down: pgdb is trying to connect, attempts left: %d, error: %s", attempts, err.Error())
		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	err = m.Down()
	defer func() { _, _ = m.Close() }()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: down error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate down: no changes")
		return nil
	}

	log.Printf("Migrate down: success")
	return nil
}
