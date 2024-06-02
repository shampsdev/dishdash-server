package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"

	"dishdash.ru/cmd/server/config"
	httpGateway "dishdash.ru/internal/gateways/http"
	"dishdash.ru/internal/repository/pg"
	"dishdash.ru/internal/usecase"
)

func main() {
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
		httpGateway.WithPort(config.C.Server.Port),
		httpGateway.WithAllowOrigin(config.C.Server.AllowOrigin),
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
