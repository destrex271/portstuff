package main

import (
	"context"
	"log"
	"net/http"
	"os"
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
		DBConnString:          getEnv("DB_CONN_STRING", "postgres://postgres:postgres@localhost:5432/trips"),
		KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
		KafkaTopic:            getEnv("KAFKA_TOPIC", "trip-requests"),
		KafkaGroupID:          getEnv("KAFKA_GROUP_ID", "trip-service-group"),
		HTTPPort:              getEnv("HTTP_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	cfg := Load()

	ctx := context.Background()

	store, err := NewStore(ctx, cfg.DBConnString)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	pub, err := NewPublisher(cfg.KafkaBootstrapServers, cfg.KafkaTopic)
	if err != nil {
		log.Fatalf("Failed to create publisher: %v", err)
	}

	cons, err := NewConsumer(cfg.KafkaBootstrapServers, cfg.KafkaTopic, cfg.KafkaGroupID)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	svc := NewService(store, pub)

	log.Println("Starting Kafka consumer")
	if err := cons.Start(ctx, svc.ProcessTripRequest); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	handler := NewHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/trips", handler.CreateTripHandler)
	mux.HandleFunc("/health", handler.HealthCheckHandler)

	log.Printf("Trip service is running on %s", cfg.HTTPPort)
	if err := http.ListenAndServe(cfg.HTTPPort, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
