package transport

import (
	"encoding/json"
	stdlog "log"
	"net/http"

	"github.com/ahmetsah/industrial-historian/go-services/audit/internal/core"
	"github.com/ahmetsah/industrial-historian/go-services/audit/internal/repository"
)

type HttpHandler struct {
	repo   repository.Repository
	hasher core.Hasher
}

func NewHttpHandler(repo repository.Repository, hasher core.Hasher) *HttpHandler {
	return &HttpHandler{repo: repo, hasher: hasher}
}

func (h *HttpHandler) Verify(w http.ResponseWriter, r *http.Request) {
	var prevHash string = "0000000000000000000000000000000000000000000000000000000000000000"
	var brokenID string
	valid := true


	err := h.repo.IterateLogs(r.Context(), func(log *core.LogEntry) error {
		if !valid {
			return nil 
		}
		
		if log.PrevHash != prevHash {
			valid = false
			brokenID = log.ID
			return nil 
		}

		calculated := h.hasher.Hash(prevHash, log)
		if log.CurrHash != calculated {
			valid = false
			brokenID = log.ID
			// Debug: stdlog the mismatch
			stdlog.Printf("Hash mismatch at %s: expected=%s, actual=%s, timestamp=%v", 
				log.ID, calculated, log.CurrHash, log.Timestamp)
			return nil
		}

		prevHash = log.CurrHash
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"valid": valid,
	}
	if !valid {
		response["broken_id"] = brokenID
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
