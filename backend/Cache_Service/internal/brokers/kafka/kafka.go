package kafka

import (
	"context"
	"errors"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader   *kafka.Reader
	log      *slog.Logger
	maxPulls int32
	jobs     chan kafka.Message
}

func NewConsumer(log *slog.Logger, brokers []string, topic string, maxPulls int32) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			MinBytes:       1,
			MaxBytes:       10e6,
			StartOffset:    kafka.LastOffset,
			CommitInterval: 0,
		}),
		log:      log,
		maxPulls: maxPulls,
		jobs:     make(chan kafka.Message, 10),
	}
}

func (c *Consumer) StartWorkerPull(ctx context.Context) error {
	const op = "Cache_Service.internal.brokers.kafka.StartWorkerPull"

	log := c.log.With(
		slog.String("op", op),
	)

	for id := range c.maxPulls {
		log.Info("Worker started", slog.Int("ID", int(id)))

		select {
		case <-ctx.Done():
			log.Info("Worker was stopped from context")
			return nil
		case _, ok := <-c.jobs:
			if !ok {
				log.Error(op, slog.String("error", "Channel read error"))
				return errors.New("channel read error")
			}
			//TODO: we need to add new method to serialise this message and call redis method
		}
	}
	return nil
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
