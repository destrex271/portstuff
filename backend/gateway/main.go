package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
)

var (
	httpAddr = EnvString("HTTP_ADDR", ":8080")

	serviceRoutes = map[string]string{
		"/auth": "http://auth-service:8081",
	}
)

func main() {
	mux := http.NewServeMux()

	for path, serviceURL := range serviceRoutes {
		targetURL, err := url.Parse(serviceURL)
		if err != nil {
			log.Printf("[ERROR]: Could not parse URL for %s: %v", path, err)
			continue
		}

		// Configure reverse proxy for the service
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		mux.Handle(path+"/", http.StripPrefix(path, proxy))
	}

	// Health check for gateway
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API Gateway is healthy")
	})

	log.Printf("API Gateway is running on %s", httpAddr)
	log.Fatal(http.ListenAndServe(httpAddr, mux))
}

func EnvString(key, fallback string) string {
	if val, ok := syscall.Getenv(key); ok {
		return val
	}
	return fallback
}
