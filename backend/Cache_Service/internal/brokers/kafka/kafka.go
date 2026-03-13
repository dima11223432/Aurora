package kafka

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	log    *slog.Logger
	jobs   chan kafka.Message
}

func NewConsumer(log *slog.Logger, brokers []string, topic string, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        groupID,
			MinBytes:       1,
			MaxBytes:       10e6,
			StartOffset:    kafka.LastOffset,
			CommitInterval: 0,
		}),
		log:  log,
		jobs: make(chan kafka.Message, 10),
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	const op = "Cache_Service.internal.brokers.kafka.Consume"

	log := c.log.With(
		slog.String("op", op),
	)

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			c.log.Error(op, slog.Any("error", err))
			return err
		}

		select {
		case c.jobs <- msg:
		case <-ctx.Done():
			log.Info("Consumer was stopped from context")
		}
	}
}

func (c *Consumer) Close() {
	c.reader.Close()
}
