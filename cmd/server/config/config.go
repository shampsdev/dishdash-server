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
	DevMode bool `default:"false" envconfig:"DEV_MODE"`
	Server  struct {
		Port uint16 `envconfig:"HTTP_PORT" default:"8000"`
	}
	DB struct {
		User        string `envconfig:"POSTGRES_USER"`
		Password    string `envconfig:"POSTGRES_PASSWORD"`
		Host        string `envconfig:"POSTGRES_HOST"`
		Port        uint16 `envconfig:"POSTGRES_PORT"`
		Database    string `envconfig:"POSTGRES_DB"`
		AutoMigrate bool   `envconfig:"POSTGRES_AUTOMIGRATE"`
	}
	TwoGisApi struct {
		Key string `envconfig:"TWOGIS_API_KEY"`
		Url string `envconfig:"TWOGIS_API_URL"`
	}
	Defaults struct {
		PriceAvg           int `default:"500" envconfig:"DEFAULT_PRICE_AVG"`
		PriceAvgUpperDelta int `default:"100" envconfig:"DEFAULT_UPPER_DELTA_AVG"`
		PriceAvgLowerDelta int `default:"300" envconfig:"DEFAULT_LOWER_DELTA_AVG"`
		Radius             int `default:"4000" envconfig:"DEFAULT_RADIUS"`
		MinDBPlaces        int `default:"5" envconfig:"DEFAULT_MIN_DB_PLACES"`
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
	if C.TwoGisApi.Url == "" {
		C.TwoGisApi.Url = "https://catalog.api.2gis.com/3.0/items"
		log.Println("[WARNING] No two gis api url, using default one: " + C.TwoGisApi.Url)
	}
	if C.TwoGisApi.Key == "" {
		log.Println("[FATAL] TwoGisApi.ApiKey is null or not set")
	}
}

func Print() {
	if C.DevMode {
		log.Println("[INFO] Launched in DEV mode")
		data, _ := json.MarshalIndent(C, "", "\t")
		fmt.Println("=== CONFIG ===")
		fmt.Println(string(data))
		fmt.Println("==============")
	} else {
		log.Println("[INFO] Launched in production mode")
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
