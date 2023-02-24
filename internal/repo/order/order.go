package order

import (
	"context"
	"errors"

	"github.com/kargotech/go-testapp/internal/entity"
	"github.com/kargotech/go-testapp/internal/repo/common"
	"github.com/kargotech/gokargo/opentelemetry/strace"
	"github.com/kargotech/gokargo/serror"
	"github.com/kargotech/gokargo/unitofwork/consistency"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

func (om OrderModule) GetOrderByID(ctx context.Context, id string) (entity.Order, error) {

	var row Order

	ctx, span := strace.Start(ctx)
	defer strace.End(span)

	err := om.db.WithContext(ctx).Where("ksuid = ?", id).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = serror.NewSuperError(ErrOrderNotFound, err, serror.WithDevMessage(err.Error()))
		} else {
			err = serror.NewSuperError(ErrOrderDataLayerInternal, err, serror.WithDevMessage(err.Error()))
		}
		return entity.Order{}, err
	}

	return entity.Order{
		Ksuid:     row.Ksuid,
		ClientID:  row.ClientID,
		ProductID: row.ProductID,
		NumSales:  row.NumSales,
	}, nil
}

func (om OrderModule) CreateOrder(ctx context.Context, orderInput entity.CreateOrder, cs consistency.ConsistencyItf) (entity.Order, error) {
	_, span := strace.Start(ctx)
	defer strace.End(span)

	tx := common.DecideDBTxn(om.db, cs)

	oKsuid := ksuid.New()
	odKsuid := ksuid.New()

	order := Order{
		Ksuid:     oKsuid.String(),
		ClientID:  orderInput.ClientID,
		ProductID: orderInput.ProductID,
		NumSales:  orderInput.NumSales,
	}

	err := tx.WithContext(ctx).Create(&order).Error
	if err != nil {
		return entity.Order{}, serror.NewSuperError(ErrOrderDataLayerInternal, err, serror.WithDevMessage(err.Error()))
	}

	orderDetails := OrderDetails{
		Ksuid:       odKsuid.String(),
		OrderKsuid:  oKsuid.String(),
		OrderType:   orderInput.OrderType,
		Origin:      orderInput.Origin,
		Destination: orderInput.Destination,
	}

	err = tx.WithContext(ctx).Create(&orderDetails).Error
	if err != nil {
		return entity.Order{}, serror.NewSuperError(ErrOrderDataLayerInternal, err, serror.WithDevMessage(err.Error()))
	}

	// transform into entity

	retOrder := entity.Order{
		Ksuid:       order.Ksuid,
		ClientID:    order.ClientID,
		ProductID:   order.ProductID,
		NumSales:    order.NumSales,
		OrderType:   orderDetails.OrderType,
		Origin:      orderDetails.Origin,
		Destination: orderDetails.Destination,
	}

	return retOrder, nil
}

func (om OrderModule) UpdateOrder(ctx context.Context, orderInput entity.UpdateOrder, cs consistency.ConsistencyItf) (entity.Order, error) {

	ctx, span := strace.Start(ctx)
	defer strace.End(span)

	tx := common.DecideDBTxn(om.db, cs)

	order := Order{
		Ksuid:    orderInput.Ksuid,
		NumSales: orderInput.NumSales,
	}

	err := tx.WithContext(ctx).Model(&order).Select("num_sales").Where("ksuid = ?", order.Ksuid).Updates(order).Error
	if err != nil {
		return entity.Order{}, serror.NewSuperError(ErrOrderDataLayerInternal, err, serror.WithDevMessage(err.Error()))
	}

	// transform into entity

	retOrder := entity.Order{
		Ksuid:    order.Ksuid,
		NumSales: order.NumSales,
	}
	return retOrder, nil
}
