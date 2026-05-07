package service

import (
	"context"
	"errors"
	"testing"

	"go_monolith_sample/internal/domain/common"
	med "go_monolith_sample/internal/domain/medicine"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Bước 1: Tạo một Mock Repository tuân thủ med.MedicineRepository
type MockMedicineRepo struct {
	mock.Mock
}

func (m *MockMedicineRepo) Create(ctx context.Context, medicine *med.Medicine) error {
	args := m.Called(ctx, medicine)
	return args.Error(0)
}

func (m *MockMedicineRepo) GetByID(ctx context.Context, id uint) (*med.Medicine, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*med.Medicine), args.Error(1)
}

func (m *MockMedicineRepo) Update(ctx context.Context, id uint, medicine *med.Medicine) error {
	args := m.Called(ctx, id, medicine)
	return args.Error(0)
}

func (m *MockMedicineRepo) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return fn(ctx)
}

func (m *MockMedicineRepo) GetAll(ctx context.Context, page, pageSize int, search string) ([]med.Medicine, *common.Pagination, error) {
	args := m.Called(ctx, page, pageSize, search)
	return args.Get(0).([]med.Medicine), args.Get(1).(*common.Pagination), args.Error(2)
}

func (m *MockMedicineRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMedicineRepo) GetByIDs(ctx context.Context, ids []uint) ([]med.Medicine, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]med.Medicine), args.Error(1)
}

func (m *MockMedicineRepo) GetByIDsForUpdate(ctx context.Context, ids []uint) ([]med.Medicine, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]med.Medicine), args.Error(1)
}

// Bước 2: Viết hàm Test cho UpdateMedicine
func TestUpdateMedicine_NotFound(t *testing.T) {
	// Khởi tạo Mock và Service
	mockRepo := new(MockMedicineRepo)
	service := NewMedicineService(mockRepo)
	ctx := context.Background()

	// Giả lập tình huống: Khi gọi GetByID(1) thì trả về lỗi "Không tìm thấy"
	mockRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	// Chạy hàm cần test
	input := med.UpdateMedicineInput{}
	err := service.UpdateMedicine(ctx, 1, input)

	// Kiểm tra kết quả
	assert.Error(t, err)
	assert.Equal(t, "Không tìm thấy thuốc để cập nhật", err.Error())

	// Đảm bảo là GetByID đã được gọi đúng 1 lần
	mockRepo.AssertExpectations(t)
}

func TestUpdateMedicine_Success(t *testing.T) {
	mockRepo := new(MockMedicineRepo)
	service := NewMedicineService(mockRepo)
	ctx := context.Background()

	oldMed := &med.Medicine{Base: common.Base{ID: 1}, Name: "Thuốc cũ"}
	newName := "Thuốc mới"
	input := med.UpdateMedicineInput{Name: &newName}

	// 1. Giả lập lấy được thuốc cũ
	mockRepo.On("GetByID", ctx, uint(1)).Return(oldMed, nil)
	// 2. Giả lập Update thành công
	mockRepo.On("Update", ctx, uint(1), mock.Anything).Return(nil)

	err := service.UpdateMedicine(ctx, 1, input)

	assert.NoError(t, err)
	assert.Equal(t, "Thuốc mới", oldMed.Name)
}
