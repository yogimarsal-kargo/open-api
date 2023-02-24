package gql

import (
	"context"

	"github.com/kargotech/go-testapp/internal/entity"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	OrderUC OrderUCItf
}

type OrderUCItf interface {
	GetOrderByID(ctx context.Context, id string) (entity.Order, error)
	CreateOrder(ctx context.Context, orderInput entity.CreateOrder) (entity.Order, error)
	UpdateOrder(ctx context.Context, orderInput entity.UpdateOrder) (entity.Order, error)
}
