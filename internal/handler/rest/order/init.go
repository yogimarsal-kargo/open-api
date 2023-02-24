package order

import (
	"context"

	"github.com/kargotech/go-testapp/internal/entity"
)

// nolint:revive
type OrderUCItf interface {
	GetOrderByID(ctx context.Context, id string) (entity.Order, error)
	CreateOrder(ctx context.Context, orderInput entity.CreateOrder) (entity.Order, error)
	UpdateOrder(ctx context.Context, orderInput entity.UpdateOrder) (entity.Order, error)
}

type Handler struct {
	uc OrderUCItf
}

func New(uc OrderUCItf) Handler {
	return Handler{
		uc: uc,
	}
}
