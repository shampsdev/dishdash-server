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
	DEBUG  bool `default:"false" envconfig:"DEBUG"`
	Server struct {
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
		PriceAvgUpperDelta int `default:"300" envconfig:"DEFAULT_UPPER_DELTA_AVG"`
		PriceAvgLowerDelta int `default:"300" envconfig:"DEFAULT_LOWER_DELTA_AVG"`
		Radius             int `default:"4000" envconfig:"DEFAULT_RADIUS"`
		MinDBPlaces        int `default:"5" envconfig:"DEFAULT_MIN_DB_PLACES"`
	}
	Recommendation struct {
		PriceCoeff float64 `default:"1" envconfig:"RECOMENDATION_PRICE_COEFF"`
		DistCoeff  float64 `default:"1" envconfig:"RECOMENDATION_DIST_COEFF"`
	}
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
	if C.TwoGisApi.Url == "" {
		C.TwoGisApi.Url = "https://catalog.api.2gis.com/3.0/items"
		log.Warn("No two gis api url, using default one: " + C.TwoGisApi.Url)
	}
	if C.TwoGisApi.Key == "" {
		log.Error("TwoGisApi.ApiKey is null or not set")
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
