package repository

import (
	"context"
	"go_monolith_sample/internal/domain/inventory"
	"go_monolith_sample/pkg/db"

	"gorm.io/gorm"
)

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) inventory.InventoryTransactionRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) Create(ctx context.Context, data *inventory.InventoryTransaction) error {
	// Luôn dùng GetTx để đảm bảo nó có thể chạy chung Transaction với Sales
	return db.GetTx(ctx, r.db).WithContext(ctx).Create(data).Error
}
