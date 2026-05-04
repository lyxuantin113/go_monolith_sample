package dto

type CreateOrderItemRequest struct {
	MedicineID uint `json:"medicine_id" validate:"required"`
	Quantity   int  `json:"quantity" validate:"required,gt=0"`
}

type CreateOrderRequest struct {
	CustomerID uint                     `json:"customer_id"`
	Items      []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}
