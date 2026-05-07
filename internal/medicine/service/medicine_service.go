package service

import (
	"context"
	"go_monolith_sample/internal/domain/common"
	med "go_monolith_sample/internal/domain/medicine"
	"go_monolith_sample/pkg/auth"
)

type medicineService struct {
	medicineRepo med.MedicineRepository
}

func NewMedicineService(medicineRepo med.MedicineRepository) *medicineService {
	return &medicineService{
		medicineRepo: medicineRepo,
	}
}

func (s *medicineService) CreateMedicine(ctx context.Context, medicine *med.Medicine) error {

	if err := medicine.Validate(); err != nil {
		return err
	}

	userID := auth.GetUserIDFromContext(ctx)
	medicine.CreatedBy = userID
	medicine.UpdatedBy = userID

	return s.medicineRepo.Create(ctx, medicine)
}

func (s *medicineService) UpdateMedicine(ctx context.Context, id uint, input med.UpdateMedicineInput) error {
	// Sử dụng Transaction để bọc toàn bộ quá trình Update
	return s.medicineRepo.Transaction(ctx, func(txCtx context.Context) error {

		// 1. Lấy dữ liệu hiện tại từ DB (Dùng txRepo để đảm bảo an toàn)
		medicine, err := s.medicineRepo.GetByID(txCtx, id)
		if err != nil {
			return med.ErrMedicineNotFound
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

		userID := auth.GetUserIDFromContext(ctx)
		medicine.UpdatedBy = userID

		return s.medicineRepo.Update(txCtx, id, medicine)
	})
}

func (s *medicineService) DeleteMedicine(ctx context.Context, id uint) error {
	return s.medicineRepo.Delete(ctx, id)
}

func (s *medicineService) GetMedicineByID(ctx context.Context, id uint) (*med.Medicine, error) {
	return s.medicineRepo.GetByID(ctx, id)
}

func (s *medicineService) GetAllMedicines(ctx context.Context, page, pageSize int, search string) ([]med.Medicine, *common.Pagination, error) {
	return s.medicineRepo.GetAll(ctx, page, pageSize, search)
}

func (s *medicineService) GetByIDs(ctx context.Context, ids []uint) ([]med.Medicine, error) {
	return s.medicineRepo.GetByIDs(ctx, ids)
}

func (s *medicineService) GetByIDsForUpdate(ctx context.Context, ids []uint) ([]med.Medicine, error) {
	return s.medicineRepo.GetByIDsForUpdate(ctx, ids)
}
