package domain

import (
	"context"
	"errors"

	"go_monolith_sample/internal/domain/common"
)

type Medicine struct {
	common.Base
	Name        string  `gorm:"type:varchar(100);notnull" json:"name" validate:"required"`
	Description string  `gorm:"type:text" json:"description"`
	Stock       int     `gorm:"default:0;notnull" json:"stock" validate:"min=0"`
	Price       float64 `gorm:"type:decimal(10,0);default:0;notnull" json:"price" validate:"gt=0"`
}

type MedicineRepository interface {
	Create(ctx context.Context, medicine *Medicine) error
	GetByID(ctx context.Context, id uint) (*Medicine, error)
	GetByIDs(ctx context.Context, ids []uint) ([]Medicine, error)
	GetByIDsForUpdate(ctx context.Context, ids []uint) ([]Medicine, error)
	GetAll(ctx context.Context, page, pageSize int, search string) ([]Medicine, *common.Pagination, error)
	Update(ctx context.Context, id uint, medicine *Medicine) error
	Delete(ctx context.Context, id uint) error
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type MedicineService interface {
	CreateMedicine(ctx context.Context, medicine *Medicine) error
	UpdateMedicine(ctx context.Context, id uint, input UpdateMedicineInput) error
	DeleteMedicine(ctx context.Context, id uint) error
	GetMedicineByID(ctx context.Context, id uint) (*Medicine, error)
	GetAllMedicines(ctx context.Context, page, pageSize int, search string) ([]Medicine, *common.Pagination, error)
	GetByIDs(ctx context.Context, ids []uint) ([]Medicine, error)
	GetByIDsForUpdate(ctx context.Context, ids []uint) ([]Medicine, error)
}

type UpdateMedicineInput struct {
	Name        *string
	Description *string
	Stock       *int
	Price       *float64
}

// Validate
var (
	ErrMedicineNotFound        = errors.New("không tìm thấy thuốc")
	ErrInvalidName             = errors.New("tên thuốc không được để trống")
	ErrInvalidPrice            = errors.New("giá thuốc phải lớn hơn 0")
	ErrInvalidStock            = errors.New("số lượng tồn kho không được âm")
	ErrInsufficientStock       = errors.New("số lượng tồn kho không đủ")
	ErrSomeMedicinesDoNotExist = errors.New("một số thuốc không tồn tại")
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
