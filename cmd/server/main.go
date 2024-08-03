package main

import (
	"context"
	"dishdash.ru/internal/usecase"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"

	"dishdash.ru/cmd/server/config"
	server "dishdash.ru/internal/gateways"
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
		log.Fatalf("can't create new database pool")
	}
	defer pool.Close()

	s := server.NewServer(usecase.Setup(pool))
	if err := s.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("error during server shutdown: %v", err)
	}
}
