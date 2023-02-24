package order

import (
	"context"
	"errors"

	"github.com/kargotech/go-testapp/internal/entity"
	repo "github.com/kargotech/go-testapp/internal/repo/order"
	"github.com/kargotech/go-testapp/pkg/validator"
	"github.com/kargotech/gokargo/opentelemetry/strace"
	"github.com/kargotech/gokargo/serror"
	"github.com/kargotech/gokargo/unitofwork"
	"github.com/kargotech/gokargo/unitofwork/consistency"
)

func (uc Usecase) GetOrderByID(ctx context.Context, id string) (entity.Order, error) {
	ctx, span := strace.Start(ctx)
	defer strace.End(span)

	order, err := uc.om.GetOrderByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrOrderNotFound):
			err = serror.NewSuperError(ErrOrderNotFound, err)
		default:
			err = serror.NewSuperError(ErrOrderInternal, err)
		}
		strace.RecordError(span, err)
		return order, err
	}
	return order, nil
}

func (uc Usecase) CreateOrder(ctx context.Context, orderInput entity.CreateOrder) (entity.Order, error) {
	ctx, span := strace.Start(ctx)
	defer strace.End(span)

	var newOrder entity.Order

	// always begin with input validation, only then proceed with
	// DB validation
	err := validator.Validator.Struct(orderInput)
	if err != nil {
		err = serror.NewSuperError(ErrOrderValidationNotPassed, err)
		strace.RecordError(span, err)
		return entity.Order{}, err
	}

	err = uc.uow.Do(ctx, &unitofwork.AuditEventInput{
		EventName: "Order Created",
		Metadata:  orderInput,
	}, func(cs consistency.ConsistencyItf) error {
		newOrder, err = uc.om.CreateOrder(ctx, orderInput, cs)
		if err != nil {
			err = serror.NewSuperError(ErrOrderInternal, err)
			strace.RecordError(span, err)
			return err
		}
		return nil
	})
	if err != nil {
		return entity.Order{}, err
	}

	return newOrder, nil
}

func (uc Usecase) UpdateOrder(ctx context.Context, orderInput entity.UpdateOrder) (entity.Order, error) {
	var updatedOrder entity.Order

	ctx, span := strace.Start(ctx)
	defer strace.End(span)

	err := validator.Validator.Struct(orderInput)
	if err != nil {
		err = serror.NewSuperError(ErrOrderValidationNotPassed, err)
		strace.RecordError(span, err)
		return entity.Order{}, err
	}

	err = uc.uow.Do(ctx, &unitofwork.AuditEventInput{
		EventName: "Order Updated",
		Metadata:  orderInput,
	}, func(cs consistency.ConsistencyItf) error {
		updatedOrder, err = uc.om.UpdateOrder(ctx, orderInput, cs)
		if err != nil {
			err = serror.NewSuperError(ErrOrderInternal, err)
			strace.RecordError(span, err)
			return err
		}
		return nil
	})
	if err != nil {
		return entity.Order{}, err
	}

	return updatedOrder, nil
}
