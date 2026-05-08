package repository

import (
	"context"
	"go_monolith_sample/internal/domain/purchase"
	"go_monolith_sample/pkg/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type purchaseOrderRepository struct {
	db *gorm.DB
}

func NewPurchaseOrderRepository(db *gorm.DB) purchase.PurchaseOrderRepository {
	return &purchaseOrderRepository{db: db}
}

func (r *purchaseOrderRepository) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx := db.GetTx(ctx, nil)
	if tx != nil {
		return fn(ctx)
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(db.InjectTx(ctx, tx))
	})
}

func (r *purchaseOrderRepository) CreatePurchaseOrder(ctx context.Context, order *purchase.PurchaseOrder) error {
	return db.GetTx(ctx, r.db).WithContext(ctx).Create(order).Error
}

func (r *purchaseOrderRepository) CreatePurchaseOrderItem(ctx context.Context, item *purchase.PurchaseOrderItem) error {
	return db.GetTx(ctx, r.db).WithContext(ctx).Create(item).Error
}

func (r *purchaseOrderRepository) GetItemsByOrderID(ctx context.Context, orderID uint) ([]purchase.PurchaseOrderItem, error) {
	var items []purchase.PurchaseOrderItem
	err := db.GetTx(ctx, r.db).WithContext(ctx).Where("purchase_order_id = ?", orderID).Find(&items).Error
	return items, err
}

func (r *purchaseOrderRepository) GetItemsByOrderIDForUpdate(ctx context.Context, orderID uint) ([]purchase.PurchaseOrderItem, error) {
	var items []purchase.PurchaseOrderItem
	err := db.GetTx(ctx, r.db).WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("purchase_order_id = ?", orderID).Find(&items).Error
	return items, err
}

func (r *purchaseOrderRepository) GetByID(ctx context.Context, id uint) (*purchase.PurchaseOrder, error) {
	var order purchase.PurchaseOrder
	err := db.GetTx(ctx, r.db).WithContext(ctx).First(&order, id).Error
	return &order, err
}

func (r *purchaseOrderRepository) GetByIDForUpdate(ctx context.Context, id uint) (*purchase.PurchaseOrder, error) {
	var order purchase.PurchaseOrder
	err := db.GetTx(ctx, r.db).WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, id).Error
	return &order, err
}

func (r *purchaseOrderRepository) Update(ctx context.Context, order *purchase.PurchaseOrder) error {
	return db.GetTx(ctx, r.db).WithContext(ctx).Save(order).Error
}

func (r *purchaseOrderRepository) Delete(ctx context.Context, id uint) error {
	return db.GetTx(ctx, r.db).WithContext(ctx).Delete(&purchase.PurchaseOrder{}, id).Error
}

func (r *purchaseOrderRepository) DeleteItemsByOrderID(ctx context.Context, orderID uint) error {
	return db.GetTx(ctx, r.db).WithContext(ctx).Where("purchase_order_id = ?", orderID).Delete(&purchase.PurchaseOrderItem{}).Error
}
