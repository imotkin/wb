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
	"github.com/imotkin/L0/internal/repo"
)

func TestGetFromCache(t *testing.T) {
	var (
		ctrl     = gomock.NewController(t)
		id       = uuid.New()
		expected = entity.Order{UID: id}
		repo     = repo.NewMockRepository(ctrl)
		cache    = cache.NewMockCache[uuid.UUID, entity.Order](ctrl)
		service  = New(logger.NewNoOp(), repo, cache)
	)

	cache.EXPECT().Get(id).Return(entity.Order{UID: id}, true)

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
		service  = New(logger.NewNoOp(), repo, cache)
	)

	cache.EXPECT().Get(id).Return(entity.Order{}, false)
	repo.EXPECT().GetOrder(gomock.Any(), id).Return(entity.Order{UID: id}, nil)
	cache.EXPECT().Set(id, entity.Order{UID: id}).Return()

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
		service = New(logger.NewNoOp(), repo, cache)
	)

	cache.EXPECT().Get(id).Return(entity.Order{}, false)
	repo.EXPECT().GetOrder(gomock.Any(), id).Return(entity.Order{}, entity.ErrOrderNotFound)

	_, err := service.Get(context.Background(), id)

	require.Equal(t, err, entity.ErrOrderNotFound)
}
