package repository

import (
	"context"
	"go_monolith_sample/internal/domain/sales"
	"go_monolith_sample/pkg/db"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) sales.OrderRepository {
	return &orderRepository{db: db}
}

func (o *orderRepository) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx := db.GetTx(ctx, nil)
	if tx != nil {
		return fn(ctx)
	}

	return o.db.Transaction(func(tx *gorm.DB) error {
		txCtx := db.InjectTx(ctx, tx)
		return fn(txCtx)
	})
}

func (o *orderRepository) Create(ctx context.Context, order *sales.Order) error {
	return o.db.WithContext(ctx).Create(order).Error
}

func (o *orderRepository) CreateItem(ctx context.Context, item *sales.OrderItem) error {
	return o.db.WithContext(ctx).Create(item).Error
}
