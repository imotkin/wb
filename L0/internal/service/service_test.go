package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/imotkin/L0/internal/cache"
	"github.com/imotkin/L0/internal/entity"
	"github.com/imotkin/L0/internal/logger"
	"github.com/imotkin/L0/internal/metrics"
	"github.com/imotkin/L0/internal/repo"
)

func TestGetFromCache(t *testing.T) {
	var (
		ctrl     = gomock.NewController(t)
		id       = uuid.New()
		expected = entity.Order{UID: id}
		repo     = repo.NewMockRepository(ctrl)
		cache    = cache.NewMockCache[uuid.UUID, entity.Order](ctrl)
		mc       = metrics.NewMockMetrics(ctrl)
		service  = New(logger.NewNoOp(), repo, cache, mc)
	)

	cache.EXPECT().Get(id).Return(entity.Order{UID: id}, true)
	mc.EXPECT().IncCacheGet()

	got, err := service.Get(context.Background(), id)

	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestGetFromRepository(t *testing.T) {
	var (
		ctrl     = gomock.NewController(t)
		id       = uuid.New()
		expected = entity.Order{UID: id}
		repo     = repo.NewMockRepository(ctrl)
		cache    = cache.NewMockCache[uuid.UUID, entity.Order](ctrl)
		mc       = metrics.NewMockMetrics(ctrl)
		service  = New(logger.NewNoOp(), repo, cache, mc)
	)

	cache.EXPECT().Get(id).Return(entity.Order{}, false)
	mc.EXPECT().IncCacheGet()

	repo.EXPECT().GetOrder(gomock.Any(), id).Return(entity.Order{UID: id}, nil)
	mc.EXPECT().IncPostgresGet()

	cache.EXPECT().Set(id, entity.Order{UID: id}).Return()
	mc.EXPECT().IncCacheSet()

	got, err := service.Get(context.Background(), id)

	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestOrderNotFound(t *testing.T) {
	var (
		ctrl    = gomock.NewController(t)
		id      = uuid.New()
		repo    = repo.NewMockRepository(ctrl)
		cache   = cache.NewMockCache[uuid.UUID, entity.Order](ctrl)
		mc      = metrics.NewMockMetrics(ctrl)
		service = New(logger.NewNoOp(), repo, cache, mc)
	)

	cache.EXPECT().Get(id).Return(entity.Order{}, false)
	mc.EXPECT().IncCacheGet()

	repo.EXPECT().GetOrder(gomock.Any(), id).Return(entity.Order{}, entity.ErrOrderNotFound)
	mc.EXPECT().IncPostgresGet()

	_, err := service.Get(context.Background(), id)

	require.ErrorIs(t, err, entity.ErrOrderNotFound)
}
