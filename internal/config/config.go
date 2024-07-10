// Package config is service-wide configuration
package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	postgresDB "github.com/nmluci/da-farm-be/internal/database/postgres"
)

var conf Config

type Config struct {
	ServiceName    string
	ServiceAddress string

	RunSince     time.Time
	PostgresConf *postgresDB.PostgresConfig
}

func New() *Config {
	if err := godotenv.Load("config/.env"); err != nil {
		panic(err)
	}

	conf = Config{
		ServiceName:    os.Getenv("SVC_NAME"),
		ServiceAddress: os.Getenv("SVC_ADDRESS"),
		RunSince:       time.Now(),
		PostgresConf: &postgresDB.PostgresConfig{
			Address:  os.Getenv("POSTGRES_ADDRESS"),
			Username: os.Getenv("POSTGRES_USERNAME"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DB:       os.Getenv("POSTGRES_DB"),
		},
	}

	return &conf
}

func Get() *Config {
	return &conf
}
