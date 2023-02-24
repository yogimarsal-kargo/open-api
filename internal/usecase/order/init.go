package order

import (
	"context"

	"github.com/kargotech/go-testapp/internal/entity"
	"github.com/kargotech/gokargo/unitofwork"
	"github.com/kargotech/gokargo/unitofwork/consistency"
)

// nolint:revive
type OrderResourceItf interface {
	GetOrderByID(ctx context.Context, id string) (entity.Order, error)
	CreateOrder(ctx context.Context, orderInput entity.CreateOrder, cs consistency.ConsistencyItf) (entity.Order, error)
	UpdateOrder(ctx context.Context, orderInput entity.UpdateOrder, cs consistency.ConsistencyItf) (entity.Order, error)
}

type Usecase struct {
	om  OrderResourceItf
	uow unitofwork.UnitOfWorkItf
}

func New(om OrderResourceItf, uow unitofwork.UnitOfWorkItf) Usecase {
	return Usecase{
		om:  om,
		uow: uow,
	}
}
