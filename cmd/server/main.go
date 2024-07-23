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
	server "dishdash.ru/internal/gateways"
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

	s := server.NewServer(setupUseCases(pool))
	if err := s.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("error during server shutdown: %v", err)
	}
}

func setupUseCases(pool *pgxpool.Pool) usecase.Cases {
	pr := pg.NewPlaceRepo(pool)
	tr := pg.NewTagRepo(pool)
	lr := pg.NewLobbyRepo(pool)
	ur := pg.NewUserRepo(pool)
	sr := pg.NewSwipeRepo(pool)

	return usecase.Cases{
		Place: usecase.NewPlaceUseCase(tr, pr),
		Tag:   usecase.NewTagUseCase(tr),
		Lobby: usecase.NewLobbyUseCase(lr, ur, tr, pr, sr),
		User:  usecase.NewUserUseCase(ur),
		Swipe: usecase.NewSwipeUseCase(sr),
	}
}
