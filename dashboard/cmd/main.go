package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"

	"dishdash.ru/pkg/usecase"

	"github.com/jackc/pgx/v5/pgxpool"

	"dashboard.dishdash.ru/cmd/config"
	server "dashboard.dishdash.ru/pkg/gateways/http"
	log "github.com/sirupsen/logrus"
)

// @title           DishDash Dashboard
// @version         1.0
// @description     Manage places
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Token

var envFile = flag.String("env-file", ".env", "Environment file")

func main() {
	flag.Parse()
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
