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

func (r *inventoryRepository) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx := db.GetTx(ctx, nil)
	if tx != nil {
		return fn(ctx)
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(db.InjectTx(ctx, tx))
	})
}

func (r *inventoryRepository) Create(ctx context.Context, data *inventory.InventoryTransaction) error {
	return db.GetTx(ctx, r.db).WithContext(ctx).Create(data).Error
}
