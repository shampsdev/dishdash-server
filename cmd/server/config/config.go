package config

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server struct {
		Port        uint16 `envconfig:"HTTP_PORT" default:"8000"`
		AllowOrigin string `envconfig:"ALLOW_ORIGIN" default:"*"`
	}
	DB struct {
		User     string `envconfig:"POSTGRES_USER"`
		Password string `envconfig:"POSTGRES_PASSWORD"`
		Host     string `envconfig:"POSTGRES_HOST"`
		Port     uint16 `envconfig:"POSTGRES_PORT"`
		Database string `envconfig:"POSTGRES_DB"`
	}
}

var C Config

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("[INFO] no .env file, parsed exported variables")
	}
	err = envconfig.Process("", &C)
	if err != nil {
		log.Fatalf("can't parse config: %s", err)
	}

	printConfig(C)
}

func printConfig(c Config) {
	data, _ := json.MarshalIndent(c, "", "\t")
	fmt.Println("=== CONFIG ===")
	fmt.Println(string(data))
	fmt.Println("==============")
}

func (c Config) DBUrl() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		C.DB.User,
		C.DB.Password,
		C.DB.Host,
		C.DB.Port,
		C.DB.Database,
	)
}

func (c Config) PGXConfig() *pgxpool.Config {
	pgxConfig, err := pgxpool.ParseConfig(c.DBUrl())
	if err != nil {
		log.Fatalf("can't parse pgxpool config: %s", err)
	}
	return pgxConfig
}
