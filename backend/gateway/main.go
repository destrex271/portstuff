package main

import (
	"log"
	"net/http"

	commons "github.com/destrex271/commons"
	_ "github.com/joho/godotenv/autoload"
)

var (
	httpAddr = commons.EnvString("HTTP_ADDR", ":8080")
)

func main() {
	mux := http.NewServeMux()
	handler := NewHandler()
	handler.RegisterRoutes(mux)

	log.Println("[INFO]: Starting HTTP Gateway Server at ", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start http server")
	}
}
