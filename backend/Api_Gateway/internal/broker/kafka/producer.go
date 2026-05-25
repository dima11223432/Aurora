// Package kafka provides a Kafka message producer for publishing events.
package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// Producer publishes messages to a Kafka topic.
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka Producer for the given brokers and topic.
func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

// Publish sends a message with the given key and value to the configured Kafka topic.
func (p *Producer) Publish(ctx context.Context, key string, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
	})

}
