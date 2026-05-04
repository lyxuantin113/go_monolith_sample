package service

import (
	"context"
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

func (s *medicineService) CreateMedicine(ctx context.Context, medicine *domain.Medicine) error {

	if err := medicine.Validate(); err != nil {
		return err
	}
	return s.medicineRepo.Create(ctx, medicine)
}

func (s *medicineService) UpdateMedicine(ctx context.Context, id uint, input domain.UpdateMedicineInput) error {
	// Sử dụng Transaction để bọc toàn bộ quá trình Update
	return s.medicineRepo.Transaction(func(txRepo domain.MedicineRepository) error {

		// 1. Lấy dữ liệu hiện tại từ DB (Dùng txRepo để đảm bảo an toàn)
		medicine, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return domain.ErrMedicineNotFound
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

		if err := medicine.Validate(); err != nil {
			return err
		}

		return txRepo.Update(ctx, id, medicine)
	})
}

func (s *medicineService) DeleteMedicine(ctx context.Context, id uint) error {
	return s.medicineRepo.Delete(ctx, id)
}

func (s *medicineService) GetMedicineByID(ctx context.Context, id uint) (*domain.Medicine, error) {
	return s.medicineRepo.GetByID(ctx, id)
}

func (s *medicineService) GetAllMedicines(ctx context.Context, page, pageSize int) ([]domain.Medicine, *domain.Pagination, error) {
	return s.medicineRepo.GetAll(ctx, page, pageSize)
}
