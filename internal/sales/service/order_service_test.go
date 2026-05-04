package service_test

import (
	"context"
	"testing"

	common "go_monolith_sample/internal/domain/common"
	inventory "go_monolith_sample/internal/domain/inventory"
	med "go_monolith_sample/internal/domain/medicine"
	medicine "go_monolith_sample/internal/domain/medicine"
	sales "go_monolith_sample/internal/domain/sales"
	service "go_monolith_sample/internal/sales/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- MOCKS ---
type MockOrderRepo struct{ mock.Mock }

func (m *MockOrderRepo) Create(ctx context.Context, order *sales.Order) error {
	return m.Called(ctx, order).Error(0)
}
func (m *MockOrderRepo) CreateItem(ctx context.Context, item *sales.OrderItem) error {
	return m.Called(ctx, item).Error(0)
}
func (m *MockOrderRepo) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type MockMedService struct{ mock.Mock }

func (m *MockMedService) GetByIDs(ctx context.Context, ids []uint) ([]med.Medicine, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]med.Medicine), args.Error(1)
}
func (m *MockMedService) GetByIDsForUpdate(ctx context.Context, ids []uint) ([]med.Medicine, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]med.Medicine), args.Error(1)
}
func (m *MockMedService) UpdateMedicine(ctx context.Context, id uint, input med.UpdateMedicineInput) error {
	return m.Called(ctx, id, input).Error(0)
}

// Giả lập các hàm khác của MedService...
func (m *MockMedService) CreateMedicine(ctx context.Context, med *med.Medicine) error {
	return nil
}
func (m *MockMedService) DeleteMedicine(ctx context.Context, id uint) error { return nil }
func (m *MockMedService) GetMedicineByID(ctx context.Context, id uint) (*med.Medicine, error) {
	return nil, nil
}
func (m *MockMedService) GetAllMedicines(ctx context.Context, p, ps int, s string) ([]med.Medicine, *common.Pagination, error) {
	return nil, nil, nil
}

type MockInvRepo struct{ mock.Mock }

func (m *MockInvRepo) Create(ctx context.Context, data *inventory.InventoryTransaction) error {
	return m.Called(ctx, data).Error(0)
}

// --- TESTS ---
func TestCreateOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("Happy Path - Success", func(t *testing.T) {
		mockOrderRepo := new(MockOrderRepo)
		mockMedService := new(MockMedService)
		mockInvRepo := new(MockInvRepo)
		s := service.NewOrderService(mockOrderRepo, mockMedService, mockInvRepo)

		input := sales.CreateOrderInput{
			CustomerID: 1,
			Items: []struct {
				MedicineID uint `json:"medicine_id" validate:"required"`
				Quantity   int  `json:"quantity" validate:"required,gt=0"`
			}{{MedicineID: 101, Quantity: 2}},
		}

		medicines := []med.Medicine{{Base: common.Base{ID: 101}, Name: "Panadol", Description: "", Stock: 10, Price: 5000}}

		// Expectation
		mockMedService.On("GetByIDs", ctx, []uint{101}).Return(medicines, nil)
		mockMedService.On("GetByIDsForUpdate", ctx, []uint{101}).Return(medicines, nil)
		mockMedService.On("UpdateMedicine", ctx, uint(101), mock.Anything).Return(nil)
		mockOrderRepo.On("Create", ctx, mock.Anything).Return(nil)
		mockOrderRepo.On("CreateItem", ctx, mock.Anything).Return(nil)
		mockInvRepo.On("Create", ctx, mock.Anything).Return(nil)

		order, err := s.CreateOrder(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, 10000.0, order.TotalPrice)
		mockMedService.AssertExpectations(t)
	})

	t.Run("Edge Case - Out of Stock", func(t *testing.T) {
		mockOrderRepo := new(MockOrderRepo)
		mockMedService := new(MockMedService)
		mockInvRepo := new(MockInvRepo)
		s := service.NewOrderService(mockOrderRepo, mockMedService, mockInvRepo)

		input := sales.CreateOrderInput{
			Items: []struct {
				MedicineID uint `json:"medicine_id" validate:"required"`
				Quantity   int  `json:"quantity" validate:"required,gt=0"`
			}{{MedicineID: 101, Quantity: 20}}, // Mua 20 trong khi chỉ có 10
		}

		medicines := []med.Medicine{{Base: common.Base{ID: 101}, Name: "Panadol", Description: "", Stock: 10, Price: 5000}}

		mockMedService.On("GetByIDs", ctx, []uint{101}).Return(medicines, nil)
		mockMedService.On("GetByIDsForUpdate", ctx, []uint{101}).Return(medicines, nil)

		_, err := s.CreateOrder(ctx, input)

		assert.ErrorIs(t, err, medicine.ErrInsufficientStock)
	})

	t.Run("Edge Case - Medicine Not Found", func(t *testing.T) {
		// ... Tương tự cho trường hợp len(medicines) != len(ids) ...
	})
}
