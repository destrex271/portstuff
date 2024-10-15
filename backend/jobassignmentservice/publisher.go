package main

// Publishes scheduled trips to trips_generated topic

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Publisher struct {
	producer *kafka.Producer
	topic    string
}

func NewPublisher(bootstrapServers, topic string) (*Publisher, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	if err != nil {
		return nil, err
	}
	return &Publisher{producer: p, topic: topic}, nil
}

func (p *Publisher) ProduceTripRequest(ctx context.Context, req TripRequest) error {
	value, err := json.Marshal(req)
	if err != nil {
		return err
	}

	deliveryChan := make(chan kafka.Event)
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, deliveryChan)

	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	return nil
}
