package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"

	"dishdash.ru/internal/repo/pg"
	"dishdash.ru/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"

	"dishdash.ru/cmd/server/config"
	server "dishdash.ru/internal/gateways"
	log "github.com/sirupsen/logrus"
)

// @title           DishDash server
// @version         2.0
// @description     Manage cards, lobbies, swipes

var envFile = flag.String("env-file", ".env", "Environment file")

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetLevel(log.DebugLevel)
	config.Load(*envFile)
	config.Print()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if config.C.DB.AutoMigrate {
		log.Info("Applying migrations")
		err := pg.MigrateUp(config.C.DBUrl())
		if err != nil {
			log.Fatalf("Can't migrate up: %s", err.Error())
		}
	} else {
		log.Info("Auto migrations is disabled")
	}
	pgConfig := config.C.PGXConfig()
	pool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		log.Fatal("can't create new database pool")
	}
	defer pool.Close()

	s := server.NewServer(usecase.Setup(pool))
	if err := s.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.WithError(err).Error("error during server shutdown")
	}
}
