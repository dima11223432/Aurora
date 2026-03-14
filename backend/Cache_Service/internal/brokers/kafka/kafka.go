package kafka

import (
	"CacheService/internal/domain/models"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type AnalyseDataProvider interface {
	SetAnalysedData(ctx context.Context, dataTitle string, data interface{}) error
}

type Consumer struct {
	reader               *kafka.Reader
	analysedDataProvider AnalyseDataProvider
	log                  *slog.Logger
	maxPulls             int32
	jobs                 chan kafka.Message
}

func NewConsumer(log *slog.Logger, brokers []string, topic string, groupID string, maxPulls int32, analyseDataProvider AnalyseDataProvider) *Consumer {
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
		analysedDataProvider: analyseDataProvider,
		log:                  log,
		maxPulls:             maxPulls,
		jobs:                 make(chan kafka.Message, 10),
	}
}

func (c *Consumer) process(ctx context.Context, msg kafka.Message) error {
	const op = "Cache_Service.internal.brokers.kafka.process"

	var analysedData models.AnalysedData

	err := json.Unmarshal(msg.Value, &analysedData)
	if err != nil {
		c.log.Error(op, slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	err = c.analysedDataProvider.SetAnalysedData(ctx, fmt.Sprintf("AnalysedData:%d", time.Now().Unix()), analysedData)
	if err != nil {
		c.log.Error(op, slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
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
					err := c.process(ctx, msg)
					if err != nil {
						log.Error(op, slog.Any("error", err))
						return
					}
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
