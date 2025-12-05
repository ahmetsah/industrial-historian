package repository_test

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ahmetsah/industrial-historian/go-services/audit/internal/core"
	"github.com/ahmetsah/industrial-historian/go-services/audit/internal/repository"
)

func TestPostgresRepository_Integration(t *testing.T) {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		t.Skip("Skipping integration test: DB_URL not set")
	}

	ctx := context.Background()
	repo, err := repository.NewPostgresRepository(ctx, dbUrl)
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}
	defer repo.Close()

	hasher := core.NewSHA256Hasher()

	// Test Concurrent Writes
	var wg sync.WaitGroup
	count := 10
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func(id int) {
			defer wg.Done()
			log := &core.LogEntry{
				Timestamp: time.Now(),
				Actor:     "tester",
				Action:    "concurrent_test",
				Details:   json.RawMessage(`{}`),
			}
			if err := repo.AppendLog(ctx, log, hasher); err != nil {
				t.Errorf("Failed to append log %d: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	// Verify Chain
	// Note: This verifies the WHOLE table. If table is large, this is slow.
	// But for test DB it's fine.
	// Also we need to handle the genesis hash logic which is implicit in IterateLogs if we start from beginning.
	// The first log in DB should have 00...00 as prevHash.
	
	var prevHash string = "0000000000000000000000000000000000000000000000000000000000000000"
	err = repo.IterateLogs(ctx, func(log *core.LogEntry) error {
		if log.PrevHash != prevHash {
			// If this is the very first log ever, it should match.
			// If we are appending to existing DB, we might start verification from middle?
			// IterateLogs starts from ORDER BY timestamp ASC.
			// So it should be the first log.
			t.Errorf("Chain broken at %s: prev %s != expected %s", log.ID, log.PrevHash, prevHash)
		}
		calculated := hasher.Hash(prevHash, log)
		if log.CurrHash != calculated {
			t.Errorf("Hash mismatch at %s", log.ID)
		}
		prevHash = log.CurrHash
		return nil
	})
	if err != nil {
		t.Errorf("IterateLogs failed: %v", err)
	}
}
