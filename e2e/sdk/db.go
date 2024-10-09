package sdk

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"

	"dishdash.ru/cmd/server/config"
	"dishdash.ru/internal/repo/pg"
)

func SetupDB() error {
	err := CleanDB()
	if err != nil {
		return err
	}

	err = pg.MigrateUp(config.C.DBUrl())
	if err != nil {
		return err
	}

	pgConfig := config.C.PGXConfig()
	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return err
	}
	sqlQuery, err := os.ReadFile("../default.sql")
	if err != nil {
		return err
	}
	_, err = pool.Exec(context.Background(), string(sqlQuery))
	if err != nil {
		return err
	}

	log.Info("Successfully setup database")
	return nil
}

func CleanDB() error {
	err := pg.MigrateDown(config.C.DBUrl())
	return err
}
