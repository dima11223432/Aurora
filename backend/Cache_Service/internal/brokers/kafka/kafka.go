package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader   *kafka.Reader
	log      *slog.Logger
	maxPulls int32
	jobs     chan kafka.Message
}

func NewConsumer(log *slog.Logger, brokers []string, topic string, groupID string, maxPulls int32) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        groupID,
			MinBytes:       1,
			MaxBytes:       10e6,
			StartOffset:    kafka.FirstOffset,
			CommitInterval: 0,
		}),
		log:      log,
		maxPulls: maxPulls,
		jobs:     make(chan kafka.Message, 10),
	}
}

func (c *Consumer) StartWorkerPull(ctx context.Context) error {
	const op = "Cache_Service.internal.brokers.kafka.StartWorkerPull"

	wg := &sync.WaitGroup{}
	wg.Add(int(c.maxPulls))

	log := c.log.With(
		slog.String("op", op),
	)

	for i := 0; i < int(c.maxPulls); i++ {

		go func(id int) {
			defer wg.Done()
			log.Info("Worker started", slog.Int("ID", int(id)))
			for {
				select {
				case <-ctx.Done():
					log.Info("Worker was stopped from context")
					return
				case msg, ok := <-c.jobs:
					if !ok {
						log.Error(op, slog.String("error", "Channel read error"))
						return
					}
					//TODO: we need to add new method to serialise this message and call redis method

				}
			}
		}(i)
	}

	wg.Wait()

	return nil
}

func (c *Consumer) Consume(ctx context.Context) error {
	const op = "Cache_Service.internal.brokers.kafka.Consume"
	defer close(c.jobs)
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
