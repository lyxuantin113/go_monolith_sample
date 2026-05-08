package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	invDomain "go_monolith_sample/internal/domain/inventory"
	medicine "go_monolith_sample/internal/domain/medicine"
	"go_monolith_sample/internal/domain/report"

	"github.com/xuri/excelize/v2"
)

type reportService struct {
	medService medicine.MedicineService
	reportRepo report.ReportRepository
}

func NewReportService(medService medicine.MedicineService, reportRepo report.ReportRepository) report.ReportService {
	return &reportService{medService: medService, reportRepo: reportRepo}
}

func (s *reportService) GetInventoryLedger(ctx context.Context, medicineID uint, from, to time.Time) ([]report.InventoryLedgerRecord, error) {
	var txs []invDomain.InventoryTransaction

	// 1. Snapshot
	lastSnapshot, err := s.reportRepo.GetLastSnapshot(ctx, medicineID, from)
	openingStock := 0
	if err == nil {
		openingStock = lastSnapshot.Stock
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Nếu là lỗi khác (mất kết nối DB...) thì mới return err
		return nil, err
	}

	// Lấy tất cả giao dịch trong khoảng thời gian, sắp xếp theo thời gian tăng dần
	txs, err = s.reportRepo.GetTransactions(ctx, medicineID, from, to)

	if err != nil {
		return nil, err
	}

	currentBalance := openingStock

	var records []report.InventoryLedgerRecord

	for _, tx := range txs {
		currentBalance += tx.Quantity
		records = append(records, report.InventoryLedgerRecord{
			Time:         tx.CreatedAt,
			Type:         string(tx.Type),
			ReferenceID:  tx.ID,
			ChangeQty:    tx.Quantity,
			BalanceAfter: currentBalance,
		})
	}

	return records, nil
}

func (s *reportService) GetSIOReport(ctx context.Context, from, to time.Time) ([]report.SIOReportItem, error) {
	var results []report.SIOReportItem

	// 1. Lấy danh sách tất cả thuốc để làm khung báo cáo
	medicines, _, err := s.medService.GetAllMedicines(ctx, 0, 0, "")
	if err != nil {
		return nil, err
	}

	for _, m := range medicines {
		item := report.SIOReportItem{
			MedicineID:   m.ID,
			MedicineName: m.Name,
		}

		// 2. Lấy Tồn đầu kỳ từ Snapshot
		lastSnapshot, err := s.reportRepo.GetLastSnapshot(ctx, m.ID, from)
		if err != nil {
			return nil, err
		}

		item.OpeningStock = lastSnapshot.Stock

		// 3. Tính Nhập và Xuất trong kỳ bằng SQL Aggregate
		type TxSummary struct {
			Inbound  int
			Outbound int
		}

		inbound, outbound, err := s.reportRepo.GetSIOSummary(ctx, m.ID, from, to)
		if err != nil {
			return nil, err
		}

		item.InboundQty = inbound
		item.OutboundQty = outbound

		// 4. Tính Tồn cuối kỳ
		item.ClosingStock = item.OpeningStock + item.InboundQty - item.OutboundQty

		results = append(results, item)
	}

	return results, nil
}

func (s *reportService) TakeInventorySnapshot(ctx context.Context) error {
	// 1. Lấy ngày hôm nay (00:00:00)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 2. Chạy một Transaction để đảm bảo tính nhất quán
	return s.reportRepo.Transaction(ctx, func(txCtx context.Context) error {
		// 3. Lấy tất cả thuốc và Stock hiện tại
		medicines, _, err := s.medService.GetAllMedicines(ctx, 0, 0, "")
		if err != nil {
			return err
		}

		// 4. Lưu vào bảng Snapshot
		var snapshots []report.InventorySnapshot
		for _, m := range medicines {
			snapshots = append(snapshots, report.InventorySnapshot{
				MedicineID:   m.ID,
				SnapshotDate: today,
				Stock:        m.Stock,
			})
		}

		return s.reportRepo.CreateSnapshots(ctx, snapshots)
	})
}

func (s *reportService) ExportSIOToExcel(ctx context.Context, from, to time.Time) ([]byte, error) {
	// 1. Lấy dữ liệu báo cáo
	data, err := s.GetSIOReport(ctx, from, to)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Bao Cao Xuat Nhap Ton"
	f.SetSheetName("Sheet1", sheet)

	// 2. Tạo Header
	headers := []string{"Mã Thuốc", "Tên Thuốc", "Tồn Đầu", "Nhập Trong Kỳ", "Xuất Trong Kỳ", "Tồn Cuối"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// 3. Điền dữ liệu
	for i, item := range data {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.MedicineID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.MedicineName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.OpeningStock)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.InboundQty)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.OutboundQty)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), item.ClosingStock)
	}

	// 4. Style: In đậm Header và kẻ bảng (Optional nhưng nên làm cho chuyên nghiệp)
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
	})
	f.SetCellStyle(sheet, "A1", "F1", style)

	// 5. Xuất ra bộ nhớ (Buffer)
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
