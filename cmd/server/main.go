package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"dishdash.ru/cmd/server/config"
	httpGateway "dishdash.ru/internal/gateways/http"
	"dishdash.ru/internal/repo/pg"
	"dishdash.ru/internal/usecase"
)

// @title           DishDash server
// @version         2.0
// @description     Manage cards, lobbies, swipes

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	pgConfig := config.C.PGXConfig()
	pool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		log.Fatalf("can't create new database pool")
	}
	defer pool.Close()

	s := httpGateway.NewServer(setupUseCases(pool))
	if err := s.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("error during server shutdown: %v", err)
	}
}

func setupUseCases(pool *pgxpool.Pool) usecase.Cases {
	cr := pg.NewCardRepository(pool)
	tr := pg.NewTagRepository(pool)
	lr := pg.NewLobbyRepository(pool)

	return usecase.Cases{
		Card:  usecase.NewCardUseCase(cr, tr),
		Tag:   usecase.NewTagUseCase(tr),
		Lobby: usecase.NewLobbyUseCase(lr),
	}
}
