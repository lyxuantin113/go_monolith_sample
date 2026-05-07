package repository

import (
	"context"
	"go_monolith_sample/internal/domain/sales"
	"go_monolith_sample/pkg/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	return o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(db.InjectTx(ctx, tx))
	})
}

func (o *orderRepository) GetByID(ctx context.Context, id uint) (*sales.Order, error) {
	var order sales.Order
	err := o.db.WithContext(ctx).First(&order, id).Error
	return &order, err
}

func (o *orderRepository) GetByIDForUpdate(ctx context.Context, id uint) (*sales.Order, error) {
	var order sales.Order
	err := o.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, id).Error
	return &order, err
}

func (o *orderRepository) GetItemsByOrderID(ctx context.Context, orderID uint) ([]sales.OrderItem, error) {
	var items []sales.OrderItem
	err := o.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&items).Error
	return items, err
}

func (o *orderRepository) GetItemsByOrderIDForUpdate(ctx context.Context, orderID uint) ([]sales.OrderItem, error) {
	var items []sales.OrderItem
	err := o.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_id = ?", orderID).Find(&items).Error
	return items, err
}

func (o *orderRepository) Update(ctx context.Context, order *sales.Order) error {
	return o.db.WithContext(ctx).Save(order).Error
}

func (o *orderRepository) Create(ctx context.Context, order *sales.Order) error {
	return o.db.WithContext(ctx).Create(order).Error
}

func (o *orderRepository) CreateItem(ctx context.Context, item *sales.OrderItem) error {
	return o.db.WithContext(ctx).Create(item).Error
}
