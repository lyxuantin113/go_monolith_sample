package domain

import (
	"context"
	"errors"
	"time"
)

type Pagination struct {
	Page     int   `json:"page" default:"1"`
	PageSize int   `json:"page_size" default:"12"`
	Total    int64 `json:"total"`
}

type Base struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}

type Medicine struct {
	Base
	Name        string  `gorm:"type:varchar(100);notnull" json:"name" validate:"required"`
	Description string  `gorm:"type:text" json:"description"`
	Stock       int     `gorm:"default:0;notnull" json:"stock" validate:"min=0"`
	Price       float64 `gorm:"type:decimal(10,0);default:0;notnull" json:"price" validate:"gt=0"`
}

type MedicineRepository interface {
	Create(ctx context.Context, medicine *Medicine) error
	GetByID(ctx context.Context, id uint) (*Medicine, error)
	GetAll(ctx context.Context, page, pageSize int) ([]Medicine, *Pagination, error)
	Update(ctx context.Context, id uint, medicine *Medicine) error
	Delete(ctx context.Context, id uint) error
	Transaction(fn func(repo MedicineRepository) error) error
}

type MedicineService interface {
	CreateMedicine(ctx context.Context, medicine *Medicine) error
	UpdateMedicine(ctx context.Context, id uint, input UpdateMedicineInput) error
	DeleteMedicine(ctx context.Context, id uint) error
	GetMedicineByID(ctx context.Context, id uint) (*Medicine, error)
	GetAllMedicines(ctx context.Context, page, pageSize int) ([]Medicine, *Pagination, error)
}

type UpdateMedicineInput struct {
	Name        *string
	Description *string
	Stock       *int
	Price       *float64
}

// Validate
var (
	ErrMedicineNotFound = errors.New("không tìm thấy thuốc")
	ErrInvalidName      = errors.New("tên thuốc không được để trống")
	ErrInvalidPrice     = errors.New("giá thuốc phải lớn hơn 0")
	ErrInvalidStock     = errors.New("số lượng tồn kho không được âm")
)

func (m *Medicine) Validate() error {
	if m.Name == "" {
		return ErrInvalidName
	}
	if m.Price <= 0 {
		return ErrInvalidPrice
	}
	if m.Stock < 0 {
		return ErrInvalidStock
	}
	return nil
}
