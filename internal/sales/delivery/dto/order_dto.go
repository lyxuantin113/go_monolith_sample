package dto

type CreateOrderItemRequest struct {
	MedicineID uint `json:"medicine_id" binding:"required" validate:"required"`
	Quantity   int  `json:"quantity" binding:"required" validate:"required,gt=0"`
}

type CreateOrderRequest struct {
	CustomerID uint                     `json:"customer_id"`
	Items      []CreateOrderItemRequest `json:"items" binding:"required" validate:"required,min=1,dive"`
}
