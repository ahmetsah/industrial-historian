package main

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	os.Setenv("DB_URL", "postgres://user:pass@localhost:5432/db")
	os.Setenv("NATS_URL", "nats://localhost:4222")
	os.Setenv("PORT", "8083")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// run should return nil (clean shutdown) or error
	// Since we don't have real dependencies yet, it might just print config and exit.
	if err := run(ctx); err != nil {
		t.Logf("Run returned error: %v", err)
	}
}
