package repository

import (
	"context"
	"go_monolith_sample/internal/domain"

	"gorm.io/gorm"
)

type medicineRepository struct {
	db *gorm.DB
}

func NewMedicineRepository(db *gorm.DB) *medicineRepository {
	return &medicineRepository{
		db: db,
	}
}

func (r *medicineRepository) WithTransaction(tx *gorm.DB) domain.MedicineRepository {
	return &medicineRepository{db: tx}
}

func (r *medicineRepository) Create(ctx context.Context, data *domain.Medicine) error {
	err := r.db.WithContext(ctx).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *medicineRepository) Update(ctx context.Context, id uint, medicine *domain.Medicine) error {
	medicine.ID = id

	err := r.db.WithContext(ctx).Save(medicine).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *medicineRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&domain.Medicine{}, "id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *medicineRepository) Transaction(fn func(domain.MedicineRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(r.WithTransaction(tx))
	})
}

func (r *medicineRepository) GetByID(ctx context.Context, id uint) (*domain.Medicine, error) {
	var medicine domain.Medicine
	err := r.db.WithContext(ctx).First(&medicine, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &medicine, nil
}

func (r *medicineRepository) GetAll(ctx context.Context, page, pageSize int) ([]domain.Medicine, *domain.Pagination, error) {
	var medicines []domain.Medicine
	var total int64

	err := r.db.WithContext(ctx).Model(&domain.Medicine{}).Count(&total).Error
	if err != nil {
		return nil, nil, err
	}

	// Tính toán offset cho phân trang
	offset := (page - 1) * pageSize

	err = r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&medicines).Error
	if err != nil {
		return nil, nil, err
	}

	pagination := &domain.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}

	return medicines, pagination, nil
}
