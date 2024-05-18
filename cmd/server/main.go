package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	httpGateway "dishdash.ru/internal/gateways/http"

	"github.com/joho/godotenv"
)

const (
	defaultPort = uint16(8000)
)

var config Config

type Config struct {
	port uint16
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	loadConfig()

	r := httpGateway.NewServer(
		httpGateway.WithPort(config.port),
	)
	if err := r.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("error during server shutdown: %v", err)
	}
}

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := defaultPort

	if portEnv, ok := os.LookupEnv("HTTP_PORT"); ok {
		port64, err := strconv.ParseInt(portEnv, 10, 16)
		if err != nil {
			log.Fatalf("Can't parse port: %s", portEnv)
		}
		port = uint16(port64)
	}

	config = Config{port: port}
}
