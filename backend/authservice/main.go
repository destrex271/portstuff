package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	commons "github.com/destrex271/commons"
	"github.com/gorilla/mux"
)

func main() {
	port, err := strconv.Atoi(commons.EnvString("port", "8081"))
	if err != nil {
		panic(err)
	}
	store, err := NewStore(commons.EnvString("connection_string", "postgres://postgres:postgres@db:5432/user"))
	if err != nil {
		panic(err)
	}
	svc := NewService(store)
	handler := NewHandler(svc)
	router := mux.NewRouter()

	// Create a subrouter for /auth
	router.HandleFunc("/register", handler.CreateUserHandler).Methods("POST")
	router.HandleFunc("/login", handler.LoginHandler).Methods("POST")
	router.HandleFunc("/health", handler.HealthCheckHandler).Methods("GET")

	// Add a catch-all route for /auth
	log.Println("Starting server on port ", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
