package healthcheck

import (
	"context"
	"sync"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imotkin/L0/internal/broker"
	"github.com/imotkin/L0/internal/metrics"
	"github.com/imotkin/L0/internal/repo/postgres"
)

func Run[T validation.Validatable](
	ctx context.Context,
	interval time.Duration,
	kafka *broker.Subscriber[T],
	pg *postgres.Postgres,
	m metrics.Metrics,
) {
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(interval):
			wg.Go(func() {
				ctx, cancel := context.WithTimeout(ctx, interval/2)
				defer cancel()

				err := kafka.Ping(ctx)
				if err != nil {
					m.SetKafkaStatus(0)
					return
				}

				m.SetKafkaStatus(1)
			})

			wg.Go(func() {
				ctx, cancel := context.WithTimeout(ctx, interval/2)
				defer cancel()

				err := pg.Ping(ctx)
				if err != nil {
					m.SetPostgresStatus(0)
					return
				}

				m.SetPostgresStatus(1)
			})

			wg.Wait()
		}
	}
}
