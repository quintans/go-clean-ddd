package infra

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	WebConfig
	DbConfig
}

type WebConfig struct {
	Port string `env:"PORT" envDefault:":8080"`
}

type DbConfig struct {
	DbName     string `env:"DB_NAME" envDefault:":postgres"`
	DbHost     string `env:"DB_HOST" envDefault:":localhost"`
	DbPort     int    `env:"DB_PORT" envDefault:"5432"`
	DbUser     string `env:"DB_USER" envDefault:"postgres"`
	DbPassword string `env:"DB_PASSWORD" envDefault:"secret"`
}

func LoadEnvVars() Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		if err != nil {
			log.Fatal(err)
		}
	}
	return cfg
}
