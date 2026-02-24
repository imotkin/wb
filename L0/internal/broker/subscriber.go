package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/imotkin/L0/internal/logger"
	"github.com/imotkin/L0/internal/metrics"
)

type Subscriber[T validation.Validatable] struct {
	r      *kgo.Client
	dlq    *kgo.Client
	values chan T
	log    logger.Logger
}

func NewConsumer[T validation.Validatable](log logger.Logger, cfg *Config) (*Subscriber[T], error) {
	reader, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Endpoint()),
		kgo.ConsumeTopics(cfg.Topic),
		kgo.ConsumerGroup(cfg.GroupID),
		kgo.DisableAutoCommit(),
	)
	if err != nil {
		return nil, fmt.Errorf("create kafka reader: %w", err)
	}

	writer, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Endpoint()),
		kgo.AllowAutoTopicCreation(),
		kgo.DefaultProduceTopic(cfg.TopicDLQ),
	)
	if err != nil {
		return nil, fmt.Errorf("create kafka dlq writer: %w", err)
	}

	return &Subscriber[T]{
		r:      reader,
		dlq:    writer,
		values: make(chan T, 10),
		log:    log.With("source", "kafka-subscriber"),
	}, nil
}

func (c *Subscriber[T]) Subscribe(ctx context.Context) <-chan T {
	c.log.Info("subscriber was started")
	go c.processMessages(ctx)
	return c.values
}

func (c *Subscriber[T]) processMessages(ctx context.Context) {
	defer c.Close()

	for {
		fetches := c.r.PollFetches(ctx)
		errs := fetches.Errors()
		if len(errs) > 0 {
			if ctx.Err() != nil {
				return
			}

			c.log.Error(errs[0].Err, "failed to fetch messages")
			continue
		}

		iter := fetches.RecordIter()

		for !iter.Done() {
			record := iter.Next()

			c.log.Info(
				"subscriber got message",
				slog.String("uid", string(record.Key)),
				slog.Int("bytes", len(record.Value)),
			)

			var value T
			err := json.Unmarshal(record.Value, &value)
			if err != nil {
				c.log.Error(err, "failed to to decode json")
				c.sendDLQ(ctx, record)
				continue
			}

			err = value.Validate()
			if err != nil {
				c.log.Error(err, "failed to validate value")
				c.sendDLQ(ctx, record)
				continue
			}

			c.values <- value
			err = c.r.CommitRecords(ctx, record)
			if err != nil {
				c.log.Error(err, "failed to commit record")
			}
		}
	}
}

func (c *Subscriber[T]) sendDLQ(ctx context.Context, record *kgo.Record) {
	res := c.dlq.ProduceSync(ctx, &kgo.Record{
		Key:   record.Key,
		Value: record.Value,
	})

	if err := res.FirstErr(); err != nil {
		c.log.Error(err, "failed to publish dlq message")
	}

	if err := c.r.CommitRecords(ctx, record); err != nil {
		c.log.Error(err, "failed to commit invalid message")
	}

	c.log.Info("message was sent to dlq")

	metrics.IncFailed()
}

func (c *Subscriber[T]) Close() {
	c.r.Close()
	c.dlq.Close()
	close(c.values)
	c.log.Info("subscriber was stopped")
}
