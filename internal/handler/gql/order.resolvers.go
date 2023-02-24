package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/kargotech/go-testapp/gen/graph/model"
	"github.com/kargotech/go-testapp/internal/entity"
)

func (r *mutationResolver) CreateOrder(ctx context.Context, input model.NewOrder) (*model.Order, error) {

	createOrderInput := entity.CreateOrder{
		ClientID:    input.ClientID,
		ProductID:   input.ProductID,
		NumSales:    input.NumSales,
		OrderType:   input.OrderType,
		Origin:      input.Origin,
		Destination: input.Origin,
	}

	result, err := r.OrderUC.CreateOrder(ctx, createOrderInput)
	if err != nil {
		return nil, err
	}

	ret := model.Order{
		ID:          result.Ksuid,
		ClientID:    result.ClientID,
		ProductID:   result.ProductID,
		NumSales:    result.NumSales,
		OrderType:   result.OrderType,
		Origin:      result.Origin,
		Destination: result.Destination,
	}

	return &ret, nil
}

func (r *mutationResolver) UpdateOrder(ctx context.Context, input model.UpdateOrder) (*model.Order, error) {

	updateOrderInput := entity.UpdateOrder{
		Ksuid:    input.ID,
		NumSales: input.NumSales,
	}

	result, err := r.OrderUC.UpdateOrder(ctx, updateOrderInput)
	if err != nil {
		return nil, err
	}

	ret := model.Order{
		ID:          result.Ksuid,
		ClientID:    result.ClientID,
		ProductID:   result.ProductID,
		NumSales:    result.NumSales,
		OrderType:   result.OrderType,
		Origin:      result.Origin,
		Destination: result.Destination,
	}

	return &ret, nil
}

func (r *queryResolver) Order(ctx context.Context, input string) (*model.Order, error) {
	result, err := r.OrderUC.GetOrderByID(ctx, input)
	if err != nil {
		return nil, err
	}

	ret := model.Order{
		ID:          result.Ksuid,
		ClientID:    result.ClientID,
		ProductID:   result.ProductID,
		NumSales:    result.NumSales,
		OrderType:   result.OrderType,
		Origin:      result.Origin,
		Destination: result.Destination,
	}

	return &ret, nil
}

func (r *queryResolver) Orders(ctx context.Context) ([]*model.Order, error) {
	panic(fmt.Errorf("not implemented"))
}
