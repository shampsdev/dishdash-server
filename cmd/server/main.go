package main

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	"os/signal"

	"dishdash.ru/cmd/server/config"
	httpGateway "dishdash.ru/internal/gateways/http"
	"dishdash.ru/internal/repository/pg"
	"dishdash.ru/internal/usecase"
)

func main() {
	config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	pgConfig := config.C.PGXConfig()
	pool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		log.Fatalf("can't create new database pool")
	}
	defer pool.Close()

	r := httpGateway.NewServer(
		setupUseCases(pool),
		httpGateway.WithPort(config.C.Port),
		httpGateway.WithAllowOrigin(config.C.AllowOrigin),
	)
	if err := r.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("error during server shutdown: %v", err)
	}
}

func setupUseCases(pool *pgxpool.Pool) usecase.Cases {

	cr := pg.NewCardRepository(pool)
	lr := pg.NewLobbyRepository(pool)
	sr := pg.NewSwipeRepository(pool)
	tr := pg.NewTagRepository(pool)

	return usecase.Cases{
		Card:  usecase.NewCard(cr, tr),
		Lobby: usecase.NewLobby(lr),
		Swipe: usecase.NewSwipe(sr),
		Tag:   usecase.NewTag(tr),
	}
}
