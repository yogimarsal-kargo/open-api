package entity

type Order struct {
	Ksuid       string `json:"ksuid"`
	ClientID    string `json:"client_id"`
	ProductID   string `json:"product_id"`
	NumSales    int    `json:"num_sales"`
	OrderType   string `json:"order_type"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}

type CreateOrder struct {
	ClientID    string `json:"client_id" validate:"required"`
	ProductID   string `json:"product_id" validate:"required"`
	NumSales    int    `json:"num_sales" validate:"required"`
	OrderType   string `json:"order_type" validate:"required"`
	Origin      string `json:"origin" validate:"required"`
	Destination string `json:"destination" validate:"required"`
}

type UpdateOrder struct {
	Ksuid    string `json:"ksuid" validate:"required"`
	NumSales int    `json:"num_sales" validate:"required"`
}
