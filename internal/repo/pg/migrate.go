package pg

import (
	"errors"
	"log"
	"time"

	"dishdash.ru/cmd/server/config"

	"github.com/golang-migrate/migrate/v4"
)

const (
	defaultAttempts = 20
	defaultTimeout  = time.Second
)

func init() {
	dbUrl := config.C.DBUrl()
	var (
		attempts = defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", dbUrl)
		if err == nil {
			break
		}

		log.Printf("Migrate: pgdb is trying to connect, attempts left: %d, error: %s", attempts, err.Error())
		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: pgdb connect error: %s", err)
	}

	err = m.Up()
	defer func() { _, _ = m.Close() }()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")
		return
	}

	log.Printf("Migrate: up success")
}
