package config

import (
	"log"
	"os"
	"strconv"

	"dishdash.ru/internal/repository/pg"
	"github.com/joho/godotenv"
)

const (
	defaultPort = uint16(8000)
)

var C = Config{}

type Config struct {
	Port uint16
	PG   pg.Config
}

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := defaultPort

	if portEnv, ok := os.LookupEnv("HTTP_PORT"); ok {
		port64, err := strconv.ParseInt(portEnv, 10, 16)
		if err != nil {
			log.Fatalf("Can't parse Port: %s", portEnv)
		}
		port = uint16(port64)
	}

	C.Port = port

	pgPort, _ := strconv.ParseInt(os.Getenv("POSTGRES_PORT"), 10, 16)

	C.PG.User = os.Getenv("POSTGRES_USER")
	C.PG.Password = os.Getenv("POSTGRES_PASSWORD")
	C.PG.Host = os.Getenv("POSTGRES_HOST")
	C.PG.Port = uint16(pgPort)
	C.PG.Database = os.Getenv("POSTGRES_DATABASE")
}
