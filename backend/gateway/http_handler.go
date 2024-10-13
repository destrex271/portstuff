package main

import "net/http"

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("POST /api/user/create", h.HandleCreateUser)
}

func (h *handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {}
