package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type Config struct {
	DBConnString          string
	KafkaBootstrapServers string
	KafkaTopic            string
	KafkaGroupID          string
	HTTPPort              string
}

func Load() Config {
	return Config{
		DBConnString:          EnvString("DB_CONN_STRING", "postgres://postgres:postgres@db:5432/trips"),
		KafkaBootstrapServers: EnvString("KAFKA_BOOTSTRAP_SERVERS", "kafka:9092"),
		KafkaTopic:            EnvString("KAFKA_TOPIC", "trip-requests"),
		KafkaGroupID:          EnvString("KAFKA_GROUP_ID", "trip-service-group"),
		HTTPPort:              EnvString("HTTP_PORT", "8082"),
	}
}

func EnvString(key, fallback string) string {
	if val, ok := syscall.Getenv(key); ok {
		return val
	}
	return fallback
}

func waitForKafka(kafkaAddr string) error {
	const maxRetries = 50
	const waitInterval = 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		// Try to connect to Kafka
		conn, err := net.Dial("tcp", kafkaAddr)
		if err == nil {
			conn.Close() // Close the connection if successful
			return nil
		}
		log.Printf("Attempting to connect to Kafka at %s: %v", kafkaAddr, err)
		time.Sleep(waitInterval)
	}

	return errors.New("could not connect to Kafka after multiple attempts")
}

func main() {
	cfg := Load()

	ctx := context.Background()

	err := waitForKafka(cfg.KafkaBootstrapServers)
	if err != nil {
		log.Fatalf("Kafka is not available: %v", err)
	}

	// Initialize store
	store, err := NewStore(ctx, cfg.DBConnString)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	// Initialize publisher
	pub, err := NewPublisher(cfg.KafkaBootstrapServers, cfg.KafkaTopic)
	if err != nil {
		log.Fatalf("Failed to create publisher: %v", err)
	}

	// Initialize consumer
	cons, err := NewConsumer(cfg.KafkaBootstrapServers, cfg.KafkaTopic, cfg.KafkaGroupID)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	// Initialize service
	svc := NewService(store, pub)

	log.Println("Starting Kafka consumer")
	if err := cons.Start(ctx, svc.ProcessTripRequest); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	handler := NewHandler(svc)
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/create", handler.CreateTripHandler).Methods("POST")
	router.HandleFunc("/health", handler.HealthCheckHandler).Methods("GET")

	port, err := strconv.Atoi(cfg.HTTPPort)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	log.Printf("Trip service is running on port %s", cfg.HTTPPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
