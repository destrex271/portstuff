package main

import (
	"context"
	"encoding/json"
	"net/http"
)

type handler struct {
	svc *service
}

func NewHandler(svc *service) *handler {
	return &handler{
		svc,
	}
}

func (handler *handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err = handler.svc.CreateUser(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func (handler *handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := handler.svc.HealthCheck(context.Background())
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"messgae": data})
}

func (handler *handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	token, err := handler.svc.LoginUser(ctx, req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}
