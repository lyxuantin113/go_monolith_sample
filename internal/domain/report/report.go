package report

import (
	"context"
	"time"

	invDomain "go_monolith_sample/internal/domain/inventory"
)

// 1. Thẻ Kho (Chi tiết 1 thuốc)
type InventoryLedgerRecord struct {
	Time         time.Time `json:"time"`
	Type         string    `json:"type"` // INBOUND, SALE, REFUND...
	ReferenceID  uint      `json:"reference_id"`
	ChangeQty    int       `json:"change_qty"`
	BalanceAfter int       `json:"balance_after"`
	Note         string    `json:"note"`
}

type InventorySnapshot struct {
	ID           uint      `gorm:"primaryKey"`
	MedicineID   uint      `gorm:"index"`
	SnapshotDate time.Time `gorm:"type:date;index"`
	Stock        int
	CreatedAt    time.Time
}

// 2. Xuất Nhập Tồn (Tổng hợp nhiều thuốc)
type SIOReportItem struct {
	MedicineID   uint   `json:"medicine_id"`
	MedicineName string `json:"medicine_name"`
	OpeningStock int    `json:"opening_stock"` // Tồn đầu kỳ
	InboundQty   int    `json:"inbound_qty"`   // Nhập trong kỳ
	OutboundQty  int    `json:"outbound_qty"`  // Xuất trong kỳ
	ClosingStock int    `json:"closing_stock"` // Tồn cuối kỳ
}

// 3. Interface Repository
type ReportRepository interface {
	GetTransactions(ctx context.Context, medicineID uint, from, to time.Time) ([]invDomain.InventoryTransaction, error)
	GetLastSnapshot(ctx context.Context, medicineID uint, before time.Time) (*InventorySnapshot, error)
	GetSIOSummary(ctx context.Context, medicineID uint, from, to time.Time) (inbound, outbound int, err error)
	CreateSnapshots(ctx context.Context, snapshots []InventorySnapshot) error
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

// 4. Interface Service
type ReportService interface {
	GetInventoryLedger(ctx context.Context, medicineID uint, from, to time.Time) ([]InventoryLedgerRecord, error)
	GetSIOReport(ctx context.Context, from, to time.Time) ([]SIOReportItem, error)
	TakeInventorySnapshot(ctx context.Context) error
	ExportSIOToExcel(ctx context.Context, from, to time.Time) ([]byte, error)
}
