package repository

import (
	"go_monolith_sample/internal/domain"

	"gorm.io/gorm"
)

type medicineRepository struct {
	db *gorm.DB
}

func NewMedicineRepository(db *gorm.DB) domain.MedicineRepository {
	return &medicineRepository{
		db: db,
	}
}

func (r *medicineRepository) WithTransaction(tx *gorm.DB) domain.MedicineRepository {
	return &medicineRepository{db: tx}
}

func (r *medicineRepository) Create(data *domain.Medicine) error {
	err := r.db.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *medicineRepository) Update(id uint, medicine *domain.Medicine) error {
	err := r.db.Model(&domain.Medicine{}).Where("id = ?", id).Updates(medicine).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *medicineRepository) Delete(id uint) error {
	err := r.db.Delete(&domain.Medicine{}, "id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *medicineRepository) GetByID(id uint) (*domain.Medicine, error) {
	var medicine domain.Medicine
	err := r.db.First(&medicine, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &medicine, nil
}

func (r *medicineRepository) GetAll() ([]domain.Medicine, error) {
	var medicines []domain.Medicine
	err := r.db.Find(&medicines).Error
	return medicines, err
}
