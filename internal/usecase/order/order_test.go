package order

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	mock_order "github.com/kargotech/go-testapp/gen/mockgen/usecase/order"
	"github.com/kargotech/go-testapp/internal/entity"
	repo "github.com/kargotech/go-testapp/internal/repo/order"
	"github.com/kargotech/gokargo/unitofwork"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_GetOrderByID(t *testing.T) {
	// Initialize Controller for mocking input
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Prepare mock resource
	mockOrderResourceItf := mock_order.NewMockOrderResourceItf(ctrl)

	type fields struct {
		om  OrderResourceItf
		uow unitofwork.UnitOfWorkItf
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func()
		want    entity.Order
		wantErr error
	}{
		// Define List of Test & Flow here
		{
			name: "ID Not found",
			fields: fields{
				om: mockOrderResourceItf,
			},
			args: args{
				ctx: context.Background(),
				id:  "anyid",
			},
			mock: func() {
				mockOrderResourceItf.EXPECT().GetOrderByID(gomock.Any(), gomock.Any()).Return(entity.Order{}, repo.ErrOrderNotFound)
			},
			wantErr: ErrOrderNotFound,
			want:    entity.Order{},
		},
		{
			name: "Success",
			fields: fields{
				om: mockOrderResourceItf,
			},
			args: args{
				ctx: context.Background(),
				id:  "anyid",
			},
			mock: func() {
				mockOrderResourceItf.EXPECT().GetOrderByID(gomock.Any(), gomock.Any()).Return(entity.Order{
					Ksuid:       "anyid",
					ClientID:    "new_client",
					ProductID:   "new_product",
					NumSales:    10,
					OrderType:   "new_type",
					Origin:      "Origin",
					Destination: "Destination",
				}, nil)
			},
			wantErr: nil,
			want: entity.Order{
				Ksuid:       "anyid",
				ClientID:    "new_client",
				ProductID:   "new_product",
				NumSales:    10,
				OrderType:   "new_type",
				Origin:      "Origin",
				Destination: "Destination",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := Usecase{
				om:  tt.fields.om,
				uow: tt.fields.uow,
			}

			// Run the mocked function defined in test case
			tt.mock()

			got, err := uc.GetOrderByID(tt.args.ctx, tt.args.id)
			if assert.Equal(t, err != nil, tt.wantErr != nil) {
				assert.ErrorIs(t, err, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUsecase_CreateOrder(t *testing.T) {
	// Initialize Controller for mocking input
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Prepare mock resource
	mockOrderResourceItf := mock_order.NewMockOrderResourceItf(ctrl)

	stubUnitOfWork := unitofwork.NewStubUnitOfWork()

	type fields struct {
		om  OrderResourceItf
		uow unitofwork.UnitOfWorkItf
	}
	type args struct {
		ctx        context.Context
		orderInput entity.CreateOrder
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func()
		want    entity.Order
		wantErr error
	}{
		// Define List of Test & Flow here
		{
			name: "Success",
			fields: fields{
				om:  mockOrderResourceItf,
				uow: stubUnitOfWork,
			},
			args: args{
				ctx: context.Background(),
				orderInput: entity.CreateOrder{
					ClientID:    "Someone",
					ProductID:   "SomeProduct",
					NumSales:    7,
					OrderType:   "APOLLO",
					Origin:      "JAWA",
					Destination: "SUMATERA",
				},
			},
			mock: func() {
				// mockUnitOfWorkItf.EXPECT().Do(gomock.Any(), gomock.Any(), gomock.Any())
				mockOrderResourceItf.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(entity.Order{Ksuid: "created"}, nil)
			},
			wantErr: nil,
			want: entity.Order{
				Ksuid: "created",
			},
		},
		{
			name: "Empty Product ID, should error Validations not passed",
			fields: fields{
				om:  mockOrderResourceItf,
				uow: stubUnitOfWork,
			},
			args: args{
				ctx: context.Background(),
				orderInput: entity.CreateOrder{
					ClientID:    "Someone",
					ProductID:   "",
					NumSales:    7,
					OrderType:   "APOLLO",
					Origin:      "JAWA",
					Destination: "SUMATERA",
				},
			},
			mock:    func() {},
			wantErr: ErrOrderValidationNotPassed,
			want:    entity.Order{},
		},
		{
			name: "Internal repo fail, should error Internal Repo Issue",
			fields: fields{
				om:  mockOrderResourceItf,
				uow: stubUnitOfWork,
			},
			args: args{
				ctx: context.Background(),
				orderInput: entity.CreateOrder{
					ClientID:    "Someone",
					ProductID:   "SomeProduct",
					NumSales:    7,
					OrderType:   "APOLLO",
					Origin:      "JAWA",
					Destination: "SUMATERA",
				},
			},
			mock: func() {
				// mockUnitOfWorkItf.EXPECT().Do(gomock.Any(), gomock.Any(), gomock.Any()).Return(serror.NewSuperError(ErrOrderInternal, nil))
				mockOrderResourceItf.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(entity.Order{}, fmt.Errorf("error insert to DB"))
			},
			wantErr: ErrOrderInternal,
			want:    entity.Order{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := Usecase{
				om:  tt.fields.om,
				uow: tt.fields.uow,
			}

			tt.mock()
			got, err := uc.CreateOrder(tt.args.ctx, tt.args.orderInput)
			if assert.Equal(t, err != nil, tt.wantErr != nil) {
				assert.ErrorIs(t, err, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUsecase_UpdateOrder(t *testing.T) {
	// Initialize Controller for mocking input
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Prepare mock resource
	mockOrderResourceItf := mock_order.NewMockOrderResourceItf(ctrl)

	stubUnitOfWork := unitofwork.NewStubUnitOfWork()

	type fields struct {
		om  OrderResourceItf
		uow unitofwork.UnitOfWorkItf
	}
	type args struct {
		ctx        context.Context
		orderInput entity.UpdateOrder
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func()
		want    entity.Order
		wantErr error
	}{
		// Define List of Test & Flow here
		{
			name: "Success",
			fields: fields{
				om:  mockOrderResourceItf,
				uow: stubUnitOfWork,
			},
			args: args{
				ctx: context.Background(),
				orderInput: entity.UpdateOrder{
					Ksuid:    "123",
					NumSales: 3,
				},
			},
			mock: func() {
				mockOrderResourceItf.EXPECT().UpdateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(entity.Order{Ksuid: "created"}, nil)
			},
			wantErr: nil,
			want: entity.Order{
				Ksuid: "created",
			},
		},
		{
			name: "Empty NumSales, should error Validations not passed",
			fields: fields{
				om:  mockOrderResourceItf,
				uow: stubUnitOfWork,
			},
			args: args{
				ctx: context.Background(),
				orderInput: entity.UpdateOrder{
					Ksuid:    "123",
					NumSales: 0,
				},
			},
			mock:    func() {},
			wantErr: ErrOrderValidationNotPassed,
			want:    entity.Order{},
		},
		{
			name: "Internal repo fail, should error Internal Repo Issue",
			fields: fields{
				om:  mockOrderResourceItf,
				uow: stubUnitOfWork,
			},
			args: args{
				ctx: context.Background(),
				orderInput: entity.UpdateOrder{
					Ksuid:    "123",
					NumSales: 3,
				},
			},
			mock: func() {
				mockOrderResourceItf.EXPECT().UpdateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(entity.Order{}, fmt.Errorf("error insert to DB"))
			},
			wantErr: ErrOrderInternal,
			want:    entity.Order{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := Usecase{
				om:  tt.fields.om,
				uow: tt.fields.uow,
			}

			tt.mock()
			got, err := uc.UpdateOrder(tt.args.ctx, tt.args.orderInput)
			if assert.Equal(t, err != nil, tt.wantErr != nil) {
				assert.ErrorIs(t, err, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
