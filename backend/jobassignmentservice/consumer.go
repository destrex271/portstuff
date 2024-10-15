package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	consumer *kafka.Consumer
}

func NewConsumer(bootstrapServers, topic, groupID string) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer: c}, nil
}

func (c *Consumer) Start(ctx context.Context, handler func(context.Context, TripRequest) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := c.consumer.ReadMessage(-1)
			if err != nil {
				log.Printf("Consumer error: %v\n", err)
				continue
			}

			var req TripRequest
			err = json.Unmarshal(msg.Value, &req)
			if err != nil {
				log.Printf("Error unmarshalling message: %v\n", err)
				continue
			}

			if err := handler(ctx, req); err != nil {
				log.Printf("Error handling message: %v\n", err)
			}
		}
	}
}
