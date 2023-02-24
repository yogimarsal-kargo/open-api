package order

import (
	"gorm.io/gorm"
)

// nolint:revive
type OrderModule struct {
	db *gorm.DB
}

func New(db *gorm.DB) OrderModule {

	// db.AutoMigrate(&model.Order{})
	// db.AutoMigrate(&model.OrderDetails{})

	return OrderModule{
		db: db,
	}

}
