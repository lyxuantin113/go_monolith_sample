package service

import (
	"context"
	"go_monolith_sample/internal/domain/common"
	"go_monolith_sample/internal/domain/inventory"
	medicine "go_monolith_sample/internal/domain/medicine"
	"go_monolith_sample/internal/domain/sales"
	"go_monolith_sample/pkg/auth"
	apperror "go_monolith_sample/pkg/error"
)

type orderService struct {
	orderRepo     sales.OrderRepository
	medService    medicine.MedicineService
	inventoryRepo inventory.InventoryTransactionRepository
}

func NewOrderService(orderRepo sales.OrderRepository, medService medicine.MedicineService, inventoryRepo inventory.InventoryTransactionRepository) sales.OrderService {
	return &orderService{
		orderRepo:     orderRepo,
		medService:    medService,
		inventoryRepo: inventoryRepo,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, input sales.CreateOrderInput) (*sales.Order, error) {
	// 1. Thu thập ID và Map số lượng (Ngoài Transaction)
	ids := make([]uint, len(input.Items))
	itemMap := make(map[uint]int)
	for i, item := range input.Items {
		ids[i] = item.MedicineID
		itemMap[item.MedicineID] = item.Quantity
	}

	// 2. Kiểm tra sơ bộ sự tồn tại (N+1)
	medicines, err := s.medService.GetByIDs(ctx, ids)
	if err != nil || len(medicines) != len(ids) {
		return nil, medicine.ErrSomeMedicinesDoNotExist
	}

	var finalOrder *sales.Order
	var inventoryLogs []*inventory.InventoryTransaction

	// 3. Mở Transaction
	err = s.orderRepo.Transaction(ctx, func(txCtx context.Context) error {
		// Chốt và Khóa dữ liệu thuốc (Race Condition)
		txMedicines, err := s.medService.GetByIDsForUpdate(txCtx, ids)
		if err != nil {
			return err
		}

		var totalAmount float64
		orderItems := make([]*sales.OrderItem, 0)

		userID := auth.GetUserIDFromContext(ctx)

		// Duyệt Stock trước
		for _, m := range txMedicines {
			qty := itemMap[m.ID]
			if m.Stock < qty {
				return medicine.ErrInsufficientStock
			}
		}

		// Thực thi nếu đủ số lượng
		for _, m := range txMedicines {
			qty := itemMap[m.ID]

			// Trừ kho
			m.Stock -= qty
			if err := s.medService.UpdateMedicine(txCtx, m.ID, medicine.UpdateMedicineInput{Stock: &m.Stock}); err != nil {
				return err
			}

			invLog := &inventory.InventoryTransaction{
				MedicineID: m.ID,
				Quantity:   qty,
				Type:       inventory.SALE,
				Price:      m.Price,
				Base: common.Base{
					CreatedBy: userID,
					// UpdatedBy: userID,
				},
			}

			// Tạm thời cho vào danh sách để lưu sau khi có OrderID
			inventoryLogs = append(inventoryLogs, invLog)

			itemTotal := float64(qty) * m.Price
			totalAmount += itemTotal

			orderItems = append(orderItems, &sales.OrderItem{
				MedicineID: m.ID,
				Quantity:   qty,
				Price:      m.Price,
				Total:      itemTotal,
				Base: common.Base{
					CreatedBy: userID,
					UpdatedBy: userID,
				},
			})

		}

		// Tạo bản ghi Đơn hàng
		order := &sales.Order{
			CustomerID: input.CustomerID,
			Status:     sales.OrderStatusCompleted,
			TotalPrice: totalAmount,
			Base: common.Base{
				CreatedBy: userID,
				UpdatedBy: userID,
			},
		}

		if err := s.orderRepo.Create(txCtx, order); err != nil {
			return err
		}

		// Lưu chi tiết đơn hàng
		for _, item := range orderItems {
			item.OrderID = order.ID
			if err := s.orderRepo.CreateItem(txCtx, item); err != nil {
				return err
			}
		}

		for _, log := range inventoryLogs {
			log.OrderID = order.ID
			if err := s.inventoryRepo.Create(txCtx, log); err != nil {
				return err
			}
		}

		finalOrder = order
		return nil
	})

	return finalOrder, err
}

func (s *orderService) RefundOrder(ctx context.Context, id uint) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil || order.Status != sales.OrderStatusCompleted {
		return apperror.BadRequest("Đơn hàng không tồn tại hoặc đã được hoàn tiền", nil)
	}

	items, err := s.orderRepo.GetItemsByOrderID(ctx, id)
	if err != nil {
		return apperror.BadRequest("Không tìm thấy chi tiết đơn hàng", nil)
	}

	ids := make([]uint, len(items))
	refundMap := make(map[uint]int)
	for i, item := range items {
		ids[i] = item.MedicineID
		refundMap[item.MedicineID] = item.Quantity
	}

	userID := auth.GetUserIDFromContext(ctx)

	err = s.orderRepo.Transaction(ctx, func(txCtx context.Context) error {

		txOrder, err := s.orderRepo.GetByIDForUpdate(txCtx, id)
		if err != nil {
			return apperror.BadRequest("Không tìm thấy đơn hàng", nil)
		}

		// Cập nhật số lượng thuốc
		txItems, err := s.medService.GetByIDsForUpdate(txCtx, ids)
		if err != nil {
			return apperror.BadRequest("Không tìm thấy chi tiết đơn hàng", nil)
		}

		for _, item := range txItems {
			refundQty := refundMap[item.ID]
			item.Stock += refundQty

			refundLog := &inventory.InventoryTransaction{
				MedicineID: item.ID,
				Quantity:   refundQty,
				Type:       inventory.REFUND,
				Price:      item.Price,
				Base: common.Base{
					CreatedBy: userID,
					UpdatedBy: userID,
				},
			}

			if err := s.medService.UpdateMedicine(txCtx, item.ID, medicine.UpdateMedicineInput{Stock: &item.Stock}); err != nil {
				return err
			}

			if err := s.inventoryRepo.Create(txCtx, refundLog); err != nil {
				return err
			}
		}

		// Status
		txOrder.Status = sales.OrderStatusRefunded
		if err := s.orderRepo.Update(txCtx, txOrder); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *orderService) DeleteOrder(ctx context.Context, orderID uint) error {
	_, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return apperror.NotFound("Không tìm thấy đơn hàng", nil)
	}

	_, err = s.orderRepo.GetItemsByOrderID(ctx, orderID)
	if err != nil {
		return apperror.NotFound("Không tìm thấy chi tiết đơn hàng", nil)
	}

	err = s.orderRepo.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.orderRepo.Delete(txCtx, orderID); err != nil {
			return err
		}

		if err := s.orderRepo.DeleteItemsByOrderID(txCtx, orderID); err != nil {
			return err
		}

		return nil
	})

	return err
}
