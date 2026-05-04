package service

import (
	"context"
	"go_monolith_sample/internal/domain/inventory"
	medicine "go_monolith_sample/internal/domain/medicine"
	"go_monolith_sample/internal/domain/sales"
)

type orderService struct {
	orderRepo     sales.OrderRepository
	medService    medicine.MedicineService
	inventoryRepo inventory.InventoryTransactionRepository
}

func NewOrderService(orderRepo sales.OrderRepository, medService medicine.MedicineService, inventoryRepo inventory.InventoryTransactionRepository) *orderService {
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

		for _, m := range txMedicines {
			qty := itemMap[m.ID]
			if m.Stock < qty {
				return medicine.ErrInsufficientStock
			}

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
			})
		}

		// Tạo bản ghi Đơn hàng
		order := &sales.Order{
			CustomerID: input.CustomerID,
			Status:     sales.OrderStatusCompleted,
			TotalPrice: totalAmount,
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
