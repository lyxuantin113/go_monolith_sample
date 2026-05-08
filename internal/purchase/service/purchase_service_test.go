package service_test

import (
	"context"
	"errors"
	"testing"

	"go_monolith_sample/internal/domain/common"
	"go_monolith_sample/internal/domain/inventory"
	medicine "go_monolith_sample/internal/domain/medicine"
	purchase "go_monolith_sample/internal/domain/purchase"
	"go_monolith_sample/internal/purchase/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- MOCKS ---
type MockPurchaseRepo struct{ mock.Mock }

func (m *MockPurchaseRepo) CreatePurchaseOrder(ctx context.Context, order *purchase.PurchaseOrder) error {
	return m.Called(ctx, order).Error(0)
}
func (m *MockPurchaseRepo) CreatePurchaseOrderItem(ctx context.Context, item *purchase.PurchaseOrderItem) error {
	return m.Called(ctx, item).Error(0)
}
func (m *MockPurchaseRepo) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}
func (m *MockPurchaseRepo) GetByID(ctx context.Context, id uint) (*purchase.PurchaseOrder, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*purchase.PurchaseOrder), args.Error(1)
}

func (m *MockPurchaseRepo) GetByIDForUpdate(ctx context.Context, id uint) (*purchase.PurchaseOrder, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*purchase.PurchaseOrder), args.Error(1)
}

func (m *MockPurchaseRepo) GetItemsByOrderID(ctx context.Context, orderID uint) ([]purchase.PurchaseOrderItem, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).([]purchase.PurchaseOrderItem), args.Error(1)
}

func (m *MockPurchaseRepo) GetItemsByOrderIDForUpdate(ctx context.Context, orderID uint) ([]purchase.PurchaseOrderItem, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).([]purchase.PurchaseOrderItem), args.Error(1)
}

func (m *MockPurchaseRepo) Update(ctx context.Context, order *purchase.PurchaseOrder) error {
	return m.Called(ctx, order).Error(0)
}

func (m *MockPurchaseRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockPurchaseRepo) DeleteItemsByOrderID(ctx context.Context, orderID uint) error {
	return m.Called(ctx, orderID).Error(0)
}

type MockMedService struct{ mock.Mock }

func (m *MockMedService) GetByIDs(ctx context.Context, ids []uint) ([]medicine.Medicine, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]medicine.Medicine), args.Error(1)
}
func (m *MockMedService) GetByIDsForUpdate(ctx context.Context, ids []uint) ([]medicine.Medicine, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]medicine.Medicine), args.Error(1)
}
func (m *MockMedService) UpdateMedicine(ctx context.Context, id uint, input medicine.UpdateMedicineInput) error {
	return m.Called(ctx, id, input).Error(0)
}
func (m *MockMedService) GetMedicineByID(ctx context.Context, id uint) (*medicine.Medicine, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*medicine.Medicine), args.Error(1)
}

// Dummy
func (m *MockMedService) CreateMedicine(ctx context.Context, med *medicine.Medicine) error {
	return nil
}
func (m *MockMedService) DeleteMedicine(ctx context.Context, id uint) error { return nil }
func (m *MockMedService) GetAllMedicines(ctx context.Context, p, ps int, s string) ([]medicine.Medicine, *common.Pagination, error) {
	return nil, nil, nil
}

type MockInvRepo struct{ mock.Mock }

func (m *MockInvRepo) Create(ctx context.Context, data *inventory.InventoryTransaction) error {
	return m.Called(ctx, data).Error(0)
}

// --- TESTS ---
func TestCreatePurchaseOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("Happy Path - Success Create", func(t *testing.T) {
		mockRepo := new(MockPurchaseRepo)
		mockMed := new(MockMedService)
		mockInv := new(MockInvRepo)
		s := service.NewPurchaseOrderService(mockRepo, mockMed, mockInv)

		input := purchase.CreatePurchaseOrderInput{
			SupplierID: 1,
			Items: []struct {
				MedicineID uint    `json:"medicine_id" binding:"required" validate:"required"`
				Quantity   int     `json:"quantity" binding:"required" validate:"required,gt=0"`
				Price      float64 `json:"price" binding:"required" validate:"required,gt=0"`
			}{{MedicineID: 10, Quantity: 100, Price: 1000}},
		}

		medicines := []medicine.Medicine{{Base: common.Base{ID: 10}, Name: "Thuốc A", Stock: 50}}

		// Bước 1: Chỉ lấy thông tin thuốc để tính giá, KHÔNG KHÓA, KHÔNG CẬP NHẬT
		mockMed.On("GetByIDs", ctx, []uint{10}).Return(medicines, nil)

		// Bước 2: Chỉ lưu Order và Item vào DB
		mockRepo.On("CreatePurchaseOrder", ctx, mock.Anything).Return(nil)
		mockRepo.On("CreatePurchaseOrderItem", ctx, mock.Anything).Return(nil)

		// QUAN TRỌNG: Lúc này CHƯA được gọi UpdateMedicine và Inventory Create
		// Nếu có gọi là logic đang bị sai (nhập kho sớm)

		po, err := s.CreatePurchaseOrder(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, po)
		assert.Equal(t, purchase.PurchaseOrderStatusPending, po.Status) // Phải là PENDING
		assert.Equal(t, 100000.0, po.TotalPrice)

		mockRepo.AssertExpectations(t)
		mockMed.AssertExpectations(t)
	})

	t.Run("Edge Case - Medicine Not Found", func(t *testing.T) {
		mockRepo := new(MockPurchaseRepo)
		mockMed := new(MockMedService)
		mockInv := new(MockInvRepo)
		s := service.NewPurchaseOrderService(mockRepo, mockMed, mockInv)

		input := purchase.CreatePurchaseOrderInput{
			Items: []struct {
				MedicineID uint    `json:"medicine_id" binding:"required" validate:"required"`
				Quantity   int     `json:"quantity" binding:"required" validate:"required,gt=0"`
				Price      float64 `json:"price" binding:"required" validate:"required,gt=0"`
			}{{MedicineID: 999, Quantity: 10, Price: 1000}},
		}

		mockMed.On("GetByIDs", ctx, []uint{999}).Return([]medicine.Medicine{}, nil)

		_, err := s.CreatePurchaseOrder(ctx, input)

		assert.ErrorIs(t, err, medicine.ErrSomeMedicinesDoNotExist)
	})

	t.Run("Happy Path - Complete Success", func(t *testing.T) {
		mockRepo := new(MockPurchaseRepo)
		mockMed := new(MockMedService)
		mockInv := new(MockInvRepo)
		s := service.NewPurchaseOrderService(mockRepo, mockMed, mockInv)

		orderID := uint(1)
		order := &purchase.PurchaseOrder{Base: common.Base{ID: orderID}, Status: purchase.PurchaseOrderStatusPending}
		items := []purchase.PurchaseOrderItem{{MedicineID: 10, Quantity: 50, Price: 1000}}
		meds := []medicine.Medicine{{Base: common.Base{ID: 10}, Stock: 100}}

		mockRepo.On("GetByIDForUpdate", mock.Anything, orderID).Return(order, nil)
		mockRepo.On("GetItemsByOrderIDForUpdate", mock.Anything, orderID).Return(items, nil)
		mockMed.On("GetByIDsForUpdate", mock.Anything, []uint{10}).Return(meds, nil)

		mockMed.On("UpdateMedicine", mock.Anything, uint(10), mock.MatchedBy(func(in medicine.UpdateMedicineInput) bool {
			return *in.Stock == 150 // 100 + 50
		})).Return(nil)
		mockInv.On("Create", mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

		err := s.CompletePurchaseOrder(ctx, orderID)

		assert.NoError(t, err)
		assert.Equal(t, purchase.PurchaseOrderStatusCompleted, order.Status)
	})

	t.Run("Edge Case - Already Completed", func(t *testing.T) {
		mockRepo := new(MockPurchaseRepo)
		s := service.NewPurchaseOrderService(mockRepo, nil, nil)

		order := &purchase.PurchaseOrder{Status: purchase.PurchaseOrderStatusCompleted}
		mockRepo.On("GetByIDForUpdate", mock.Anything, uint(1)).Return(order, nil)

		err := s.CompletePurchaseOrder(ctx, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "đã được xử lý")
	})

	t.Run("Edge Case - DB Error During Inbound Update", func(t *testing.T) {
		mockRepo := new(MockPurchaseRepo)
		mockMed := new(MockMedService)
		mockInv := new(MockInvRepo)
		s := service.NewPurchaseOrderService(mockRepo, mockMed, mockInv)

		input := purchase.CreatePurchaseOrderInput{
			Items: []struct {
				MedicineID uint    `json:"medicine_id" binding:"required" validate:"required"`
				Quantity   int     `json:"quantity" binding:"required" validate:"required,gt=0"`
				Price      float64 `json:"price" binding:"required" validate:"required,gt=0"`
			}{{MedicineID: 10, Quantity: 100, Price: 1000}},
		}

		medicines := []medicine.Medicine{{Base: common.Base{ID: 10}, Name: "Thuốc A", Stock: 50, Price: 2000}}

		mockMed.On("GetByIDs", ctx, []uint{10}).Return(medicines, nil)
		mockMed.On("GetByIDsForUpdate", ctx, []uint{10}).Return(medicines, nil)
		mockRepo.On("CreatePurchaseOrder", ctx, mock.Anything).Return(nil)

		// Giả lập lỗi khi cập nhật kho
		mockRepo.On("CreatePurchaseOrderItem", ctx, mock.Anything).Return(errors.New("db connection timeout"))
		_, err := s.CreatePurchaseOrder(ctx, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timeout")
	})
}
