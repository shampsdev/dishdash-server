package config

import (
	"encoding/json"
	"fmt"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug  bool `default:"false" envconfig:"DEBUG"`
	Server struct {
		Port uint16 `envconfig:"HTTP_PORT" default:"8000"`
	}
	DB struct {
		User     string `envconfig:"POSTGRES_USER"`
		Password string `envconfig:"POSTGRES_PASSWORD"`
		Host     string `envconfig:"POSTGRES_HOST"`
		Port     uint16 `envconfig:"POSTGRES_PORT"`
		Database string `envconfig:"POSTGRES_DB"`
	}
}

func Load(envFile string) *Config {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Info("no .env file, parsed exported variables")
	}
	c := &Config{}
	err = envconfig.Process("", c)
	if err != nil {
		log.Fatalf("can't parse config: %s", err)
	}
	return c
}

func (c *Config) Print() {
	if c.Debug {
		log.Info("Launched in debug mode")
		data, _ := json.MarshalIndent(c, "", "\t")
		fmt.Println("=== CONFIG ===")
		fmt.Println(string(data))
		fmt.Println("==============")
	} else {
		log.Info("Launched in production mode")
	}
}

func (c *Config) DBUrl() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Database,
	)
}

func (c Config) PGXConfig() *pgxpool.Config {
	pgxConfig, err := pgxpool.ParseConfig(c.DBUrl())
	if err != nil {
		log.Fatalf("can't parse pgxpool config: %s", err)
	}
	return pgxConfig
}
