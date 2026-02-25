package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"

	"github.com/imotkin/L0/internal/api/handler"
	"github.com/imotkin/L0/internal/api/router"
	"github.com/imotkin/L0/internal/api/server"
	"github.com/imotkin/L0/internal/broker"
	"github.com/imotkin/L0/internal/cache"
	"github.com/imotkin/L0/internal/config"
	"github.com/imotkin/L0/internal/entity"
	"github.com/imotkin/L0/internal/logger"
	"github.com/imotkin/L0/internal/repo/postgres"
	"github.com/imotkin/L0/internal/service"
)

func TestOrder() (key string, v any) {
	id := uuid.New()

	return id.String(), entity.Order{
		UID:         id,
		TrackNumber: uuid.NewString(),
		Entry:       "WBIL",
		Delivery: entity.Delivery{
			Name:    "Иван Иванов",
			Phone:   "+79999999999",
			Zip:     "101000",
			City:    "Москва",
			Address: "Площадь Мира, стр. 15",
			Region:  "Центральный",
			Email:   "ivanov@example.com",
		},
		Payment: entity.Payment{
			Transaction:  uuid.New(),
			RequestID:    uuid.NewString(),
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []entity.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				RID:         uuid.New(),
				Name:        "Product 1",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "ABC",
				Status:      202,
			},
			{
				ChrtID:      9934931,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				RID:         uuid.New(),
				Name:        "Product 2",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "DEF",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "sign-123",
		CustomerID:        uuid.NewString(),
		DeliveryService:   "DHL",
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		Shard:             "1",
	}
}

var configPath = flag.String("config", "config.example.yaml", "path to config file")

func Run() error {
	flag.Parse()

	cfg, err := config.Parse(*configPath)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	err = cfg.Validate()
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	lw, err := logger.ParseOutput(cfg.Logging.Output)
	if err != nil {
		return fmt.Errorf("parse logger output: %w", err)
	}

	log := logger.New(cfg.Logging.Format, cfg.Logging.Level, lw)

	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM}

	ctx, cancel := signal.NotifyContext(context.Background(), signals...)
	defer cancel()

	pg, err := postgres.New(ctx, cfg.Postgres.ConnectionURL())
	if err != nil {
		return fmt.Errorf("create postgres: %w", err)
	}

	err = pg.MigrateUp(ctx, cfg.Postgres.MigrationsDir)
	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	pub, err := broker.NewPublisher(log, cfg.Broker)
	if err != nil {
		return fmt.Errorf("create producer: %w", err)
	}

	pub.IntervalPublish(ctx, TestOrder, cfg.Broker.Interval)

	sub, err := broker.NewConsumer[entity.Order](log, cfg.Broker)
	if err != nil {
		return fmt.Errorf("create producer: %w", err)
	}

	var (
		c = cache.New[uuid.UUID, entity.Order](cfg.Cache.Size)
		s = service.New(log, pg, c)
		h = handler.New(log, s)
		r = router.New(h, cfg.Web.TemplatePath)
	)

	s.Run(ctx, sub)

	return server.New(log, cfg.Server, r).Start(ctx)
}
