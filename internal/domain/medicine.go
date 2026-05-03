package domain

import "time"

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
	Create(medicine *Medicine) error
	GetByID(id uint) (*Medicine, error)
	GetAll() ([]Medicine, error)
	Update(id uint, medicine *Medicine) error
	Delete(id uint) error
	Transaction(fn func(repo MedicineRepository) error) error
}

type MedicineService interface {
	CreateMedicine(medicine *Medicine) error
	UpdateMedicine(id uint, input UpdateMedicineInput) error
	DeleteMedicine(id uint) error
	GetMedicineByID(id uint) (*Medicine, error)
	GetAllMedicines() ([]Medicine, error)
}

type UpdateMedicineInput struct {
	Name        *string
	Description *string
	Stock       *int
	Price       *float64
}
