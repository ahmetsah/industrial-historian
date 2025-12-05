package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func main() {
	// Simulate the hash calculation
	prevHash := "0000000000000000000000000000000000000000000000000000000000000000"
	timestamp := "2025-12-04 07:32:34.396973 +0000 UTC"
	actor := "admin"
	action := "login"
	details := `{"ip": "127.0.0.1"}`
	
	// Parse timestamp
	t, _ := time.Parse("2006-01-02 15:04:05.999999 -0700 MST", timestamp)
	
	// Calculate hash as the hasher does
	data := fmt.Sprintf("%s%s%s%s%s",
		prevHash,
		t.Format(time.RFC3339Nano),
		actor,
		action,
		details,
	)
	
	hash := sha256.Sum256([]byte(data))
	result := hex.EncodeToString(hash[:])
	
	fmt.Printf("Data: %s\n", data)
	fmt.Printf("Hash: %s\n", result)
	fmt.Printf("Expected: 96108dfbc15ba92eb9bf4a3e87b4bed507c31bc03af71cc73ad5814d21dba396\n")
	fmt.Printf("Actual:   10768136ec14abc15c1318a58cc10603ad3bde474716eac491de2cc346b365b0\n")
}
