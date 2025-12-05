package config

import (
	"fmt"
	"os"
)

type Config struct {
	DbUrl   string
	NatsUrl string
	Port    string
}

func LoadConfig() (*Config, error) {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("DB_URL environment variable is required")
	}

	natsUrl := os.Getenv("NATS_URL")
	if natsUrl == "" {
		return nil, fmt.Errorf("NATS_URL environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	return &Config{
		DbUrl:   dbUrl,
		NatsUrl: natsUrl,
		Port:    port,
	}, nil
}
