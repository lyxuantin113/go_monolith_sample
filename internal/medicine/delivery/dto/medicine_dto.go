package dto

type CreateMedicineRequest struct {
	Name        string  `json:"name" binding:"required" validate:"required,min=3"`
	Price       float64 `json:"price" binding:"required" validate:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required" validate:"required,min=0"`
	Description string  `json:"description" binding:"omitempty" validate:"omitempty"`
}

type UpdateMedicineRequest struct {
	Name        *string  `json:"name" binding:"omitempty" validate:"omitempty,min=3"`
	Price       *float64 `json:"price" binding:"omitempty" validate:"omitempty,gt=0"`
	Stock       *int     `json:"stock" binding:"omitempty" validate:"omitempty,min=0"`
	Description *string  `json:"description" binding:"omitempty" validate:"omitempty"`
}
