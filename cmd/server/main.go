package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"dishdash.ru/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"

	"dishdash.ru/cmd/server/config"
	server "dishdash.ru/internal/gateways"
	log "github.com/sirupsen/logrus"
)

// @title           DishDash server
// @version         2.0
// @description     Manage cards, lobbies, swipes

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config.Print()
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

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetLevel(log.DebugLevel)
}
