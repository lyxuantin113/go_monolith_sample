package repository

import (
	"context"
	"time"

	invDomain "go_monolith_sample/internal/domain/inventory"
	"go_monolith_sample/internal/domain/report"
	"go_monolith_sample/pkg/db"

	"gorm.io/gorm"
)

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) report.ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx := db.GetTx(ctx, nil)
	if tx != nil {
		return fn(ctx)
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(db.InjectTx(ctx, tx))
	})
}

func (r *reportRepository) GetTransactions(ctx context.Context, medicineID uint, from, to time.Time) ([]invDomain.InventoryTransaction, error) {
	var txs []invDomain.InventoryTransaction
	err := r.db.WithContext(ctx).
		Where("medicine_id = ? AND created_at BETWEEN ? AND ?", medicineID, from, to).
		Order("created_at asc").
		Find(&txs).Error
	return txs, err
}

func (r *reportRepository) GetLastSnapshot(ctx context.Context, medicineID uint, before time.Time) (*report.InventorySnapshot, error) {
	var snapshot report.InventorySnapshot
	err := r.db.WithContext(ctx).
		Where("medicine_id = ? AND snapshot_date < ?", medicineID, before).
		Order("snapshot_date desc").
		First(&snapshot).Error
	return &snapshot, err
}

func (r *reportRepository) GetSIOSummary(ctx context.Context, medicineID uint, from, to time.Time) (inbound, outbound int, err error) {
	type TxSummary struct {
		Inbound  int
		Outbound int
	}

	var summary TxSummary
	err = r.db.WithContext(ctx).
		Model(&invDomain.InventoryTransaction{}).
		Select("SUM(CASE WHEN quantity > 0 THEN quantity ELSE 0 END) as inbound, SUM(CASE WHEN quantity < 0 THEN ABS(quantity) ELSE 0 END) as outbound").
		Where("medicine_id = ? AND created_at BETWEEN ? AND ?", medicineID, from, to).
		Scan(&summary).Error

	return summary.Inbound, summary.Outbound, err
}

func (r *reportRepository) CreateSnapshots(ctx context.Context, snapshots []report.InventorySnapshot) error {
	return r.db.WithContext(ctx).Create(&snapshots).Error
}
