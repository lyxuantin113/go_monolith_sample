package service

import (
	"errors"
	"testing"

	"go_monolith_sample/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Bước 1: Tạo một Mock Repository tuân thủ domain.MedicineRepository
type MockMedicineRepo struct {
	mock.Mock
}

func (m *MockMedicineRepo) Create(medicine *domain.Medicine) error {
	args := m.Called(medicine)
	return args.Error(0)
}

func (m *MockMedicineRepo) GetByID(id uint) (*domain.Medicine, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Medicine), args.Error(1)
}

func (m *MockMedicineRepo) Update(id uint, medicine *domain.Medicine) error {
	args := m.Called(id, medicine)
	return args.Error(0)
}

func (m *MockMedicineRepo) Transaction(fn func(domain.MedicineRepository) error) error {
	// Với Unit Test, ta giả lập transaction bằng cách chạy thẳng hàm fn luôn
	return fn(m)
}

// Giả lập các hàm khác cho đủ Interface...
func (m *MockMedicineRepo) GetAll() ([]domain.Medicine, error) { return nil, nil }
func (m *MockMedicineRepo) Delete(id uint) error               { return nil }

// Bước 2: Viết hàm Test cho UpdateMedicine
func TestUpdateMedicine_NotFound(t *testing.T) {
	// Khởi tạo Mock và Service
	mockRepo := new(MockMedicineRepo)
	service := NewMedicineService(mockRepo)

	// Giả lập tình huống: Khi gọi GetByID(1) thì trả về lỗi "Không tìm thấy"
	mockRepo.On("GetByID", uint(1)).Return(nil, errors.New("not found"))

	// Chạy hàm cần test
	input := domain.UpdateMedicineInput{}
	err := service.UpdateMedicine(1, input)

	// Kiểm tra kết quả
	assert.Error(t, err)
	assert.Equal(t, "Không tìm thấy thuốc để cập nhật", err.Error())

	// Đảm bảo là GetByID đã được gọi đúng 1 lần
	mockRepo.AssertExpectations(t)
}

func TestUpdateMedicine_Success(t *testing.T) {
	mockRepo := new(MockMedicineRepo)
	service := NewMedicineService(mockRepo)

	oldMed := &domain.Medicine{Base: domain.Base{ID: 1}, Name: "Thuốc cũ"}
	newName := "Thuốc mới"
	input := domain.UpdateMedicineInput{Name: &newName}

	// 1. Giả lập lấy được thuốc cũ
	mockRepo.On("GetByID", uint(1)).Return(oldMed, nil)
	// 2. Giả lập Update thành công
	mockRepo.On("Update", uint(1), mock.Anything).Return(nil)

	err := service.UpdateMedicine(1, input)

	assert.NoError(t, err)
	assert.Equal(t, "Thuốc mới", oldMed.Name)
}
