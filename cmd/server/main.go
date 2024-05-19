package main

import (
	"context"
	"errors"
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

	r := httpGateway.NewServer(
		setupUseCases(ctx),
		httpGateway.WithPort(config.C.Port),
		httpGateway.WithAllowOrigin(config.C.AllowOrigin),
	)
	if err := r.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("error during server shutdown: %v", err)
	}
}

func setupUseCases(ctx context.Context) usecase.Cases {
	db, err := pg.NewPostgresDB(ctx, pg.Config{
		User:     config.C.PG.User,
		Password: config.C.PG.Password,
		Host:     config.C.PG.Host,
		Port:     config.C.PG.Port,
		Database: config.C.PG.Database,
	})
	if err != nil {
		log.Fatalf("can't setup postgres: %s", err)
	}

	cr := pg.NewCardRepository(db)
	lr := pg.NewLobbyRepository(db)
	sr := pg.NewSwipeRepository(db)

	return usecase.Cases{
		Card:  usecase.NewCard(cr),
		Lobby: usecase.NewLobby(lr),
		Swipe: usecase.NewSwipe(sr),
	}
}
