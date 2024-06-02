package config

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"strconv"

	"dishdash.ru/internal/repository/pg"
	"github.com/joho/godotenv"
)

// TODO конфиг через https://github.com/sethvargo/go-envconfig
const (
	defaultPort = uint16(8000)
)

var C = Config{}

type Config struct {
	Port        uint16
	AllowOrigin string
	PG          pg.Config
}

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
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
	C.AllowOrigin = os.Getenv("ALLOW_ORIGIN")

	pgPort, _ := strconv.ParseInt(os.Getenv("POSTGRES_PORT"), 10, 16)

	C.PG.User = os.Getenv("POSTGRES_USER")
	C.PG.Password = os.Getenv("POSTGRES_PASSWORD")
	C.PG.Host = os.Getenv("POSTGRES_HOST")
	C.PG.Port = uint16(pgPort)
	C.PG.Database = os.Getenv("POSTGRES_DATABASE")
}

func (c Config) PGXConfig() *pgxpool.Config {
	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		C.PG.User,
		C.PG.Password,
		C.PG.Host,
		C.PG.Port,
		C.PG.Database,
	)
	pgxConfig, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		log.Fatalf("can't parse pgxpool config: %s", err)
	}
	return pgxConfig
}
