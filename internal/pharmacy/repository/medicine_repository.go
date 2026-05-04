package repository

import (
	"context"

	common "go_monolith_sample/internal/domain/common"
	med "go_monolith_sample/internal/domain/medicine"
	db "go_monolith_sample/pkg/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type medicineRepository struct {
	db *gorm.DB
}

func NewMedicineRepository(db *gorm.DB) *medicineRepository {
	return &medicineRepository{
		db: db,
	}
}

func (r *medicineRepository) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx := db.GetTx(ctx, nil)
	if tx != nil {
		return fn(ctx)
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(db.InjectTx(ctx, tx))
	})
}

func (r *medicineRepository) Create(ctx context.Context, data *med.Medicine) error {
	return db.GetTx(ctx, r.db).WithContext(ctx).Create(data).Error
}

func (r *medicineRepository) Update(ctx context.Context, id uint, medicine *med.Medicine) error {
	medicine.ID = id
	return db.GetTx(ctx, r.db).WithContext(ctx).Save(medicine).Error
}

func (r *medicineRepository) Delete(ctx context.Context, id uint) error {
	return db.GetTx(ctx, r.db).WithContext(ctx).Delete(&med.Medicine{}, "id = ?", id).Error
}

func (r *medicineRepository) GetByID(ctx context.Context, id uint) (*med.Medicine, error) {
	var medicine med.Medicine
	err := db.GetTx(ctx, r.db).WithContext(ctx).First(&medicine, "id = ?", id).Error
	return &medicine, err
}

func (r *medicineRepository) GetByIDs(ctx context.Context, ids []uint) ([]med.Medicine, error) {
	var medicines []med.Medicine
	err := db.
		GetTx(ctx, r.db).
		WithContext(ctx).
		Find(&medicines, "id IN ?", ids).
		Error
	return medicines, err
}

func (r *medicineRepository) GetByIDsForUpdate(ctx context.Context, ids []uint) ([]med.Medicine, error) {
	var medicines []med.Medicine
	// Clause Locking Strength UPDATE tương đương với SELECT ... FOR UPDATE
	err := db.GetTx(ctx, r.db).WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Find(&medicines, "id IN ?", ids).Error
	return medicines, err
}

func (r *medicineRepository) GetAll(ctx context.Context, page, pageSize int, search string) ([]med.Medicine, *common.Pagination, error) {
	var medicines []med.Medicine
	var total int64

	query := db.GetTx(ctx, r.db).WithContext(ctx).Model(&med.Medicine{})

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, nil, err
	}

	// Tính toán offset cho phân trang
	offset := (page - 1) * pageSize

	err = query.Offset(offset).Limit(pageSize).Find(&medicines).Error
	if err != nil {
		return nil, nil, err
	}

	pagination := &common.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}

	return medicines, pagination, nil
}
