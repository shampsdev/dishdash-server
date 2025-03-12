package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shampsdev/dishdash-server/crudoshlep/pkg/config"
	"github.com/shampsdev/dishdash-server/crudoshlep/pkg/gateways/rest"
	"github.com/shampsdev/dishdash-server/crudoshlep/pkg/usecase"
	log "github.com/sirupsen/logrus"
)

// @title           Crudoshlep server
// @version         1.0
// @description     Save analytics events
func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetLevel(log.DebugLevel)

	log.Info("Hello from tglinked server!")

	cfg := config.Load(".env")
	cfg.Print()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	pgConfig := cfg.PGXConfig()
	pool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		log.Fatal("can't create new database pool")
	}
	defer pool.Close()

	s := rest.NewServer(cfg, usecase.Setup(pool))
	if err := s.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.WithError(err).Error("error during server shutdown")
	}
}
