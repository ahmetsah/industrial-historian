package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("DB_URL", "postgres://user:pass@localhost:5432/db")
	os.Setenv("NATS_URL", "nats://localhost:4222")
	os.Setenv("PORT", "8083")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.DbUrl != "postgres://user:pass@localhost:5432/db" {
		t.Errorf("Expected DB_URL to be %s, got %s", "postgres://user:pass@localhost:5432/db", cfg.DbUrl)
	}
	if cfg.NatsUrl != "nats://localhost:4222" {
		t.Errorf("Expected NATS_URL to be %s, got %s", "nats://localhost:4222", cfg.NatsUrl)
	}
	if cfg.Port != "8083" {
		t.Errorf("Expected PORT to be %s, got %s", "8083", cfg.Port)
	}
}

func TestLoadConfig_MissingEnv(t *testing.T) {
	os.Unsetenv("DB_URL")
	_, err := LoadConfig()
	if err == nil {
		t.Error("Expected error when DB_URL is missing, got nil")
	}
}
