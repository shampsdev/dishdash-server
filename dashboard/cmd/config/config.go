package config

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DEBUG  bool `default:"true" envconfig:"DEBUG"`
	Server struct {
		Port uint16 `envconfig:"HTTP_PORT"    default:"8000"`
	}
	DB struct {
		User     string `envconfig:"POSTGRES_USER"`
		Password string `envconfig:"POSTGRES_PASSWORD"`
		Host     string `envconfig:"POSTGRES_HOST"`
		Port     uint16 `envconfig:"POSTGRES_PORT"`
		Database string `envconfig:"POSTGRES_DB"`
	}
	Auth struct {
		ApiToken string `envconfig:"API_TOKEN"`
	}
	Parser struct {
		URL    string `envconfig:"PARSER_URL"`
		ApiKey string `envconfig:"PARSER_API_KEY"`
	}
	S3 S3Config
}

type S3Config struct {
	AccessKeyID string `envconfig:"S3_ACCESS_KEY_ID"`
	SecretKey   string `envconfig:"S3_SECRET_KEY"`
	Region      string `envconfig:"S3_REGION"`
	Bucket      string `envconfig:"S3_BUCKET"`
	EndpointUrl string `envconfig:"S3_ENDPOINT_URL"`
}

var C Config

func Load(envFile string) {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Info("no .env file, parsed exported variables")
	}
	err = envconfig.Process("", &C)
	if err != nil {
		log.Fatalf("can't parse config: %s", err)
	}
}

func Print() {
	if C.DEBUG {
		log.Info("Launched in DEV mode")
		data, _ := json.MarshalIndent(C, "", "\t")
		fmt.Println("=== CONFIG ===")
		fmt.Println(string(data))
		fmt.Println("==============")
	} else {
		log.Info("Launched in production mode")
	}
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
