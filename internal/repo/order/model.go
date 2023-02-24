package order

import (
	"time"

	auditresource "github.com/kargotech/gokargo/audit/resource"
	"gorm.io/gorm"
)

type Order struct {
	Ksuid     string `json:"ksuid"`
	ClientID  string `json:"client_id"`
	ProductID string `json:"product_id"`
	NumSales  int    `json:"num_sales"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// nolint:revive
type OrderDetails struct {
	Ksuid       string `json:"ksuid"`
	OrderKsuid  string `json:"order_ksuid"`
	OrderType   string `json:"order_type"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

func (o *Order) AfterCreate(tx *gorm.DB) error {
	_ = auditresource.AutoCreate(tx, o.Ksuid)

	return nil
}

func (o *Order) BeforeUpdate(tx *gorm.DB) error {
	_ = auditresource.AutoUpdate(tx, o.Ksuid)

	return nil
}

func (o *OrderDetails) AfterCreate(tx *gorm.DB) error {
	_ = auditresource.AutoCreate(tx, o.Ksuid)

	return nil
}

func (o *OrderDetails) BeforeUpdate(tx *gorm.DB) error {
	_ = auditresource.AutoUpdate(tx, o.Ksuid)

	return nil
}
