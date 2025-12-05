package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("DB_URL", "postgres://user:pass@localhost:5432/audit")
	os.Setenv("NATS_URL", "nats://localhost:4222")
	defer os.Unsetenv("DB_URL")
	defer os.Unsetenv("NATS_URL")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.DbUrl != "postgres://user:pass@localhost:5432/audit" {
		t.Errorf("Expected DB_URL to be 'postgres://user:pass@localhost:5432/audit', got '%s'", cfg.DbUrl)
	}
	if cfg.NatsUrl != "nats://localhost:4222" {
		t.Errorf("Expected NATS_URL to be 'nats://localhost:4222', got '%s'", cfg.NatsUrl)
	}
}
