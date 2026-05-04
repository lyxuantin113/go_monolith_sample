package inventory

import (
	"context"
	"go_monolith_sample/internal/domain/common"
)

type InventoryTransactionRepository interface {
	Create(ctx context.Context, data *InventoryTransaction) error
}

type TypeInventoryTransaction string

const (
	INBOUND TypeInventoryTransaction = "INBOUND"
	SALE    TypeInventoryTransaction = "SALE"
	REFUND  TypeInventoryTransaction = "REFUND"
)

type InventoryTransaction struct {
	common.Base
	OrderID    uint
	MedicineID uint
	Quantity   int
	Type       TypeInventoryTransaction
	Price      float64 `gorm:"type:decimal(10,0);default:0;notnull" json:"price" validate:"gt=0"`
}
