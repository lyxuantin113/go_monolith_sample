package service

import (
	"errors"
	"go_monolith_sample/internal/domain"
)

type medicineService struct {
	medicineRepo domain.MedicineRepository
}

func NewMedicineService(medicineRepo domain.MedicineRepository) domain.MedicineService {
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
	// Chuyển đổi dữ liệu (Mapping) từ Input sang Model
	medicine := &domain.Medicine{}

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

	// Thực hiện cập nhật
	return s.medicineRepo.Update(id, medicine)
}

func (s *medicineService) DeleteMedicine(id uint) error {

	if id == 0 {
		return errors.New("ID thuốc không được để trống")
	}

	return s.medicineRepo.Delete(id)
}

func (s *medicineService) GetMedicineByID(id uint) (*domain.Medicine, error) {

	if id == 0 {
		return nil, errors.New("ID thuốc không được để trống")
	}

	return s.medicineRepo.GetByID(id)
}

func (s *medicineService) GetAllMedicines() ([]domain.Medicine, error) {
	return s.medicineRepo.GetAll()
}
