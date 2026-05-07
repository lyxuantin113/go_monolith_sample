package purchase

import (
	"context"

	"go_monolith_sample/internal/domain/common"
)

type Supplier struct {
	common.Base
	Name    string `gorm:"type:varchar(100);notnull" json:"name"`
	Email   string `gorm:"type:varchar(100);unique" json:"email"`
	Phone   string `gorm:"type:varchar(20)" json:"phone"`
	Address string `gorm:"type:varchar(255)" json:"address"`
}

type PurchaseOrderStatus string

const (
	PurchaseOrderStatusPending   PurchaseOrderStatus = "PENDING"
	PurchaseOrderStatusCompleted PurchaseOrderStatus = "COMPLETED"
	PurchaseOrderStatusCancelled PurchaseOrderStatus = "CANCELLED"
)

type PurchaseOrder struct {
	common.Base
	SupplierID uint                `gorm:"notnull" json:"supplier_id"`
	Status     PurchaseOrderStatus `gorm:"type:varchar(50);notnull" json:"status" validate:"required"`
	TotalPrice float64             `gorm:"type:decimal(10,0);default:0;notnull" json:"total_price" validate:"gt=0"`
}

type PurchaseOrderItem struct {
	common.Base
	PurchaseOrderID uint    `gorm:"notnull" json:"purchase_order_id"`
	MedicineID      uint    `gorm:"notnull" json:"medicine_id"`
	Quantity        int     `gorm:"notnull" json:"quantity" validate:"gt=0"`
	Price           float64 `gorm:"type:decimal(10,0);default:0;notnull" json:"price" validate:"gt=0"`
}

type PurchaseOrderRepository interface {
	CreatePurchaseOrder(ctx context.Context, order *PurchaseOrder) error
	CreatePurchaseOrderItem(ctx context.Context, item *PurchaseOrderItem) error
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error

	GetByID(ctx context.Context, id uint) (*PurchaseOrder, error)
	GetByIDForUpdate(ctx context.Context, id uint) (*PurchaseOrder, error)
	GetItemsByOrderID(ctx context.Context, orderID uint) ([]PurchaseOrderItem, error)
	GetItemsByOrderIDForUpdate(ctx context.Context, orderID uint) ([]PurchaseOrderItem, error)
	Update(ctx context.Context, order *PurchaseOrder) error
}

type PurchaseOrderService interface {
	CreatePurchaseOrder(ctx context.Context, input CreatePurchaseOrderInput) (*PurchaseOrder, error)
	CompletePurchaseOrder(ctx context.Context, id uint) error
	CancelPurchaseOrder(ctx context.Context, id uint) error
}

type CreatePurchaseOrderInput struct {
	SupplierID uint `json:"supplier_id" binding:"required" validate:"required"`
	Items      []struct {
		MedicineID uint    `json:"medicine_id" binding:"required" validate:"required"`
		Quantity   int     `json:"quantity" binding:"required" validate:"required,gt=0"`
		Price      float64 `json:"price" binding:"required" validate:"required,gt=0"`
	} `json:"items" binding:"required" validate:"required,min=1,dive"`
}
