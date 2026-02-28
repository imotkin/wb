package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/imotkin/L0/internal/broker"
	"github.com/imotkin/L0/internal/cache"
	"github.com/imotkin/L0/internal/entity"
	"github.com/imotkin/L0/internal/logger"
	"github.com/imotkin/L0/internal/metrics"
	"github.com/imotkin/L0/internal/repo"
)

type OrderService struct {
	cache cache.Cache[uuid.UUID, entity.Order]
	repo  repo.Repository
	log   logger.Logger
	mc    metrics.Metrics
}

func New(
	log logger.Logger,
	repo repo.Repository,
	cache cache.Cache[uuid.UUID, entity.Order],
	mc metrics.Metrics,
) *OrderService {
	return &OrderService{
		cache: cache,
		repo:  repo,
		log:   log.With("source", "order-service"),
		mc:    mc,
	}
}

func (s *OrderService) Get(ctx context.Context, id uuid.UUID) (entity.Order, error) {
	order, ok := s.cache.Get(id)
	s.mc.IncCacheGet()

	if ok {
		return order, nil
	}

	order, err := s.repo.GetOrder(ctx, id)
	s.mc.IncPostgresGet()

	if err != nil {
		return entity.Order{}, fmt.Errorf("get from repository: %w", err)
	}

	s.cache.Set(id, order)
	s.mc.IncCacheSet()

	return order, nil
}

func (s *OrderService) List(ctx context.Context) ([]entity.Order, error) {
	return s.repo.List(ctx)
}

func (s *OrderService) Add(ctx context.Context, order entity.Order) (bool, error) {
	return s.repo.AddOrder(ctx, order)
}

func (s *OrderService) initCache(ctx context.Context) {
	orders, err := s.repo.List(ctx)
	if err != nil {
		s.log.Error(err, "failed to init cache from database")
		return
	}

	for _, order := range orders {
		s.cache.Set(order.UID, order)
	}

	s.log.Info("cache was inited", "size", s.cache.Len())
}

func (s *OrderService) processOrder(ctx context.Context, order entity.Order) {
	_, ok := s.cache.Get(order.UID)
	s.mc.IncCacheGet()

	if ok {
		s.log.Warn("duplicate order was sent", "uid", order.UID)
		return
	}

	inserted, err := s.Add(ctx, order)
	if err != nil {
		s.log.Error(err, "failed to add order", "uid", order.UID)
		return
	}

	if !inserted {
		s.log.Warn("duplicate order was sent", "uid", order.UID)
		return
	}

	s.log.Info("order was added", "uid", order.UID)
	s.mc.IncOrders()

	s.cache.Set(order.UID, order)
	s.mc.IncCacheSet()
}

func (s *OrderService) Run(ctx context.Context, sub *broker.Subscriber[entity.Order]) {
	s.initCache(ctx)

	orders := sub.Subscribe(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case order, ok := <-orders:
				if !ok {
					return
				}

				s.processOrder(ctx, order)
			}
		}
	}()
}
