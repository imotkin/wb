package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/imotkin/L0/internal/logger"
)

type Publisher struct {
	c   *kgo.Client
	log logger.Logger
}

func NewPublisher(log logger.Logger, cfg *Config) (*Publisher, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers([]string{cfg.Endpoint()}...),
		kgo.AllowAutoTopicCreation(),
		kgo.DefaultProduceTopic(cfg.Topic),
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		c:   client,
		log: log.With("source", "kafka-publisher"),
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, key string, value any) (int, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return 0, fmt.Errorf("encode value: %w", err)
	}

	r := p.c.ProduceSync(
		context.Background(),
		&kgo.Record{
			Key:   []byte(key),
			Value: bytes,
		},
	)

	if err := r.FirstErr(); err != nil {
		return 0, err
	}

	return len(bytes), nil
}

func (p *Publisher) IntervalPublish(ctx context.Context, fn func() (string, any), interval time.Duration) {
	p.log.Info("publisher was started", slog.Duration("interval", interval))

	go func() {
		defer p.log.Info("publisher was stopped")

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.Tick(interval):
				key, value := fn()

				size, err := p.Publish(ctx, key, value)
				if err != nil {
					p.log.Error(
						err,
						"failed to publish message",
						slog.String("uid", key),
					)
					continue
				}

				p.log.Info(
					"publisher sent message",
					slog.String("uid", key),
					slog.Int("bytes", size),
				)
			}
		}
	}()
}
