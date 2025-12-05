package core

import (
	"encoding/json"
	"time"
)

type LogEntry struct {
	ID        string          `json:"id"`
	Timestamp time.Time       `json:"timestamp"`
	Actor     string          `json:"actor"`
	Action    string          `json:"action"`
	Details   json.RawMessage `json:"details"`
	PrevHash  string          `json:"prev_hash"`
	CurrHash  string          `json:"curr_hash"`
}
