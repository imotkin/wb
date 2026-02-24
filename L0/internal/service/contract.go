package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/imotkin/L0/internal/entity"
)

type Service interface {
	Add(ctx context.Context, order entity.Order) (bool, error)
	Get(ctx context.Context, id uuid.UUID) (entity.Order, error)
}
