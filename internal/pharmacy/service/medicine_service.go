package service

import (
	"errors"
	"go_monolith_sample/internal/domain"
)

type medicineService struct {
	medicineRepo domain.MedicineRepository
}

func NewMedicineService(medicineRepo domain.MedicineRepository) *medicineService {
	return &medicineService{
		medicineRepo: medicineRepo,
	}
}

func (s *medicineService) CreateMedicine(medicine *domain.Medicine) error {

	if medicine.Name == "" {
		return errors.New("Tên thuốc không được để trống")
	}

	if medicine.Stock < 0 {
		return errors.New("Số lượng tồn kho không được âm")
	}

	if medicine.Price <= 0 {
		return errors.New("Giá thuốc phải lớn hơn 0")
	}

	return s.medicineRepo.Create(medicine)
}

func (s *medicineService) UpdateMedicine(id uint, input domain.UpdateMedicineInput) error {
	// Sử dụng Transaction để bọc toàn bộ quá trình Update
	return s.medicineRepo.Transaction(func(txRepo domain.MedicineRepository) error {

		// 1. Lấy dữ liệu hiện tại từ DB (Dùng txRepo để đảm bảo an toàn)
		medicine, err := txRepo.GetByID(id)
		if err != nil {
			return errors.New("Không tìm thấy thuốc để cập nhật")
		}

		// 2. Mapping và chuẩn bị dữ liệu mới (Patch)
		if input.Name != nil {
			medicine.Name = *input.Name
		}
		if input.Price != nil {
			medicine.Price = *input.Price
		}
		if input.Stock != nil {
			medicine.Stock = *input.Stock
		}
		if input.Description != nil {
			medicine.Description = *input.Description
		}

		// 3. Thực hiện cập nhật xuống DB (Vẫn dùng txRepo)
		return txRepo.Update(id, medicine)
	})
}

func (s *medicineService) DeleteMedicine(id uint) error {
	return s.medicineRepo.Delete(id)
}

func (s *medicineService) GetMedicineByID(id uint) (*domain.Medicine, error) {
	return s.medicineRepo.GetByID(id)
}

func (s *medicineService) GetAllMedicines() ([]domain.Medicine, error) {
	return s.medicineRepo.GetAll()
}
