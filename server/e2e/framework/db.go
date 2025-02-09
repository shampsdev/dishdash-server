package framework

import (
	"context"
	"fmt"
	"os"

	"dishdash.ru/cmd/server/config"
	"dishdash.ru/pkg/repo/pg"
)

func (fw *Framework) SetupDB() error {
	err := pg.MigrateDown(fw.Cfg.DBUrl())
	if err != nil {
		return err
	}

	err = pg.MigrateUp(config.C.DBUrl())
	if err != nil {
		return fmt.Errorf("failed to migrate up: %w", err)
	}

	sqlQuery, err := os.ReadFile("../default.sql")
	if err != nil {
		return fmt.Errorf("failed to read sql query: %w", err)
	}

	_, err = fw.DB.Exec(context.Background(), string(sqlQuery))
	if err != nil {
		return fmt.Errorf("failed to execute sql query: %w", err)
	}

	fw.Log.Infof("Setup database")

	return nil
}
