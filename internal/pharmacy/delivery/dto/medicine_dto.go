package dto

type CreateMedicineRequest struct {
	Name        string  `json:"name" validate:"required,min=3"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"required,min=0"`
	Description string  `json:"description"`
}

type UpdateMedicineRequest struct {
	Name        *string  `json:"name" validate:"omitempty,min=3"`
	Price       *float64 `json:"price" validate:"omitempty,gt=0"`
	Stock       *int     `json:"stock" validate:"omitempty,min=0"`
	Description *string  `json:"description" validate:"omitempty"`
}
