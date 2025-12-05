package core

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSHA256Hasher_Hash(t *testing.T) {
	hasher := NewSHA256Hasher()
	
	timestamp, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	log := &LogEntry{
		Timestamp: timestamp,
		Actor:     "admin",
		Action:    "login",
		Details:   json.RawMessage(`{"ip":"127.0.0.1"}`),
	}
	prevHash := "0000000000000000000000000000000000000000000000000000000000000000"

	hash := hasher.Hash(prevHash, log)

	if len(hash) != 64 {
		t.Errorf("Expected hash length 64, got %d", len(hash))
	}

	// Verify determinism
	hash2 := hasher.Hash(prevHash, log)
	if hash != hash2 {
		t.Errorf("Hash should be deterministic")
	}

	// Verify sensitivity
	log.Actor = "user"
	hash3 := hasher.Hash(prevHash, log)
	if hash == hash3 {
		t.Errorf("Hash should change when content changes")
	}
}
