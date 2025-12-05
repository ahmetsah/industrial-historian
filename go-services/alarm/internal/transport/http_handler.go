package transport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ahmetsah/industrial-historian/go-services/alarm/internal/core"
)

type HttpHandler struct {
	service *core.AlarmService
}

func NewHttpHandler(service *core.AlarmService) *HttpHandler {
	return &HttpHandler{service: service}
}

func (h *HttpHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/alarms/{id}/ack", h.handleAck)
	mux.HandleFunc("POST /api/v1/alarms/{id}/shelve", h.handleShelve)
	mux.HandleFunc("GET /api/v1/alarms/active", h.handleListActive)
	mux.HandleFunc("POST /api/v1/alarms/definitions", h.handleCreateDefinition)
}

func (h *HttpHandler) handleAck(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid alarm ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Acknowledge(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"acknowledged"}`))
}

func (h *HttpHandler) handleShelve(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid alarm ID", http.StatusBadRequest)
		return
	}

	var req struct {
		DurationSeconds int `json:"duration_seconds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	duration := time.Duration(req.DurationSeconds) * time.Second
	if err := h.service.Shelve(id, duration); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"shelved"}`))
}

func (h *HttpHandler) handleListActive(w http.ResponseWriter, r *http.Request) {
	alarms := h.service.GetActiveAlarms()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alarms)
}

func (h *HttpHandler) handleCreateDefinition(w http.ResponseWriter, r *http.Request) {
	var def core.AlarmDefinition
	if err := json.NewDecoder(r.Body).Decode(&def); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateDefinition(&def); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(def)
}
