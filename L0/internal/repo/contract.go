package repo

import (
	"context"

	"github.com/google/uuid"

	"github.com/imotkin/L0/internal/entity"
)

type Repository interface {
	AddOrder(ctx context.Context, order entity.Order) (bool, error)
	GetOrder(ctx context.Context, id uuid.UUID) (entity.Order, error)
	List(ctx context.Context) ([]entity.Order, error)
}
