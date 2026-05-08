package sales

import (
	"context"
	"go_monolith_sample/internal/domain/common"
)

type OrderStatus string

const (
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusRefunded  OrderStatus = "refunded"
)

type Order struct {
	common.Base
	CustomerID uint        `gorm:"notnull" json:"customer_id"`
	TotalPrice float64     `gorm:"type:decimal(10,0);default:0;notnull" json:"total_price" validate:"gt=0"`
	Status     OrderStatus `gorm:"type:varchar(100);notnull" json:"status" validate:"required"`
}

type OrderItem struct {
	common.Base
	OrderID    uint    `gorm:"notnull" json:"order_id"`
	MedicineID uint    `gorm:"notnull" json:"medicine_id"`
	Quantity   int     `gorm:"notnull" json:"quantity" validate:"gt=0"`
	Price      float64 `gorm:"type:decimal(10,0);default:0;notnull" json:"price" validate:"gt=0"`
	Total      float64 `gorm:"type:decimal(10,0);default:0;notnull" json:"total" validate:"gt=0"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	CreateItem(ctx context.Context, item *OrderItem) error
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error

	GetByID(ctx context.Context, id uint) (*Order, error)
	GetByIDForUpdate(ctx context.Context, id uint) (*Order, error)
	GetItemsByOrderID(ctx context.Context, orderID uint) ([]OrderItem, error)
	GetItemsByOrderIDForUpdate(ctx context.Context, orderID uint) ([]OrderItem, error)
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id uint) error
	DeleteItemsByOrderID(ctx context.Context, orderID uint) error
}

type OrderService interface {
	CreateOrder(ctx context.Context, input CreateOrderInput) (*Order, error)
	RefundOrder(ctx context.Context, orderID uint) error
	DeleteOrder(ctx context.Context, orderID uint) error
}

type CreateOrderInput struct {
	CustomerID uint `json:"customer_id"`
	Items      []struct {
		MedicineID uint `json:"medicine_id" validate:"required"`
		Quantity   int  `json:"quantity" validate:"required,gt=0"`
	} `json:"items" validate:"required,min=1,dive"`
}
