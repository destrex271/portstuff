package main

import (
	"context"
	"encoding/json"
	"net/http"
)

type Handler struct {
	svc JobAssignmentService
}

func NewHandler(svc JobAssignmentService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateTripHandler(w http.ResponseWriter, r *http.Request) {
	var req TripRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate trip request
	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	if err := h.svc.RequestTrip(ctx, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Trip requested successfully"})
}

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Service is healthy"})
}
