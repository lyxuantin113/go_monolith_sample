package service

import (
	"context"
	"go_monolith_sample/internal/domain/inventory"
	medicine "go_monolith_sample/internal/domain/medicine"
	purchase "go_monolith_sample/internal/domain/purchase"
	apperror "go_monolith_sample/pkg/error"
)

type purchaseOrderService struct {
	repo          purchase.PurchaseOrderRepository
	medService    medicine.MedicineService
	inventoryRepo inventory.InventoryTransactionRepository
}

func NewPurchaseOrderService(
	repo purchase.PurchaseOrderRepository,
	medService medicine.MedicineService,
	inventoryRepo inventory.InventoryTransactionRepository,
) purchase.PurchaseOrderService {
	return &purchaseOrderService{
		repo:          repo,
		medService:    medService,
		inventoryRepo: inventoryRepo,
	}
}

func (s *purchaseOrderService) CreatePurchaseOrder(ctx context.Context, input purchase.CreatePurchaseOrderInput) (*purchase.PurchaseOrder, error) {

	// 1. IDs và Map items
	ids := make([]uint, len(input.Items))
	itemMap := make(map[uint]struct {
		qty   int
		price float64
	})

	for i, item := range input.Items {
		ids[i] = item.MedicineID
		itemMap[item.MedicineID] = struct {
			qty   int
			price float64
		}{item.Quantity, item.Price}
	}

	// 2. N+1
	medicines, err := s.medService.GetByIDs(ctx, ids)
	if err != nil || len(medicines) != len(ids) {
		return nil, medicine.ErrSomeMedicinesDoNotExist
	}

	// 3. Transaction
	var finalOrder *purchase.PurchaseOrder

	err = s.repo.Transaction(ctx, func(txCtx context.Context) error {

		// Tạo Order
		var totalAmount float64
		for _, item := range input.Items {
			totalAmount += float64(item.Quantity) * item.Price
		}

		order := &purchase.PurchaseOrder{
			SupplierID: input.SupplierID,
			Status:     purchase.PurchaseOrderStatusPending,
			TotalPrice: totalAmount,
		}

		if err := s.repo.CreatePurchaseOrder(txCtx, order); err != nil {
			return err
		}

		// Xử lý từng Item
		for _, m := range medicines {
			inputInfo := itemMap[m.ID]

			// Lưu chi tiết đơn nhập
			if err := s.repo.CreatePurchaseOrderItem(txCtx, &purchase.PurchaseOrderItem{
				PurchaseOrderID: order.ID,
				MedicineID:      m.ID,
				Quantity:        inputInfo.qty,
				Price:           inputInfo.price,
			}); err != nil {
				return err
			}
		}

		finalOrder = order
		return nil
	})

	return finalOrder, err
}

func (s *purchaseOrderService) CompletePurchaseOrder(ctx context.Context, id uint) error {
	return s.repo.Transaction(ctx, func(txCtx context.Context) error {
		// 1. Lấy và KHÓA đơn hàng để tránh bị Cancel/Complete cùng lúc từ máy khác
		order, err := s.repo.GetByIDForUpdate(txCtx, id)
		if err != nil || order.Status != purchase.PurchaseOrderStatusPending {
			return apperror.BadRequest("Đơn hàng không hợp lệ hoặc đã được xử lý", err)
		}

		// 2. Lấy và KHÓA danh sách món hàng
		txItems, err := s.repo.GetItemsByOrderIDForUpdate(txCtx, id)
		if err != nil || len(txItems) == 0 {
			return apperror.NotFound("Không tìm thấy món hàng nào trong đơn hàng", err)
		}

		// Gom ID thuốc để lấy hàng loạt
		medIDs := make([]uint, len(txItems))
		for i, item := range txItems {
			medIDs[i] = item.MedicineID
		}

		// 3. Khóa các loại thuốc trong DB và tạo Map để truy xuất nhanh
		medicines, err := s.medService.GetByIDsForUpdate(txCtx, medIDs)
		if err != nil {
			return err
		}

		medMap := make(map[uint]medicine.Medicine)
		for _, m := range medicines {
			medMap[m.ID] = m
		}

		// 4. Cập nhật tồn kho
		for _, item := range txItems {
			med, exists := medMap[item.MedicineID]
			if !exists {
				return apperror.Internal("Dữ liệu thuốc không đồng bộ", nil)
			}

			newStock := med.Stock + item.Quantity
			if err := s.medService.UpdateMedicine(txCtx, med.ID, medicine.UpdateMedicineInput{Stock: &newStock}); err != nil {
				return err
			}

			// Ghi nhật ký kho (INBOUND)
			if err := s.inventoryRepo.Create(txCtx, &inventory.InventoryTransaction{
				MedicineID: item.MedicineID,
				OrderID:    order.ID,
				Quantity:   item.Quantity,
				Type:       inventory.INBOUND,
				Price:      item.Price,
			}); err != nil {
				return err
			}
		}

		// 5. Chốt đơn
		order.Status = purchase.PurchaseOrderStatusCompleted
		return s.repo.Update(txCtx, order)
	})
}

func (s *purchaseOrderService) CancelPurchaseOrder(ctx context.Context, id uint) error {
	return s.repo.Transaction(ctx, func(txCtx context.Context) error {
		// 1. Lấy thông tin đơn hàng
		order, err := s.repo.GetByIDForUpdate(txCtx, id)
		if err != nil || order.Status != purchase.PurchaseOrderStatusPending {
			return apperror.NotFound("Không tìm thấy đơn hàng hoặc đơn hàng đã được xử lý", err)
		}

		// 2. Hủy đơn hàng
		order.Status = purchase.PurchaseOrderStatusCancelled
		return s.repo.Update(txCtx, order)
	})
}
