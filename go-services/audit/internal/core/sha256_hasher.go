package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type SHA256Hasher struct{}

func NewSHA256Hasher() *SHA256Hasher {
	return &SHA256Hasher{}
}

func (h *SHA256Hasher) Hash(prevHash string, log *LogEntry) string {
	// Truncate to Microsecond to match Postgres storage precision
	ts := log.Timestamp.Truncate(time.Microsecond)

	// Format: prevHash + timestamp(RFC3339Nano) + actor + action + details
	data := fmt.Sprintf("%s%s%s%s%s",
		prevHash,
		ts.Format(time.RFC3339Nano),
		log.Actor,
		log.Action,
		string(log.Details),
	)

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
