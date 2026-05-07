package dto

type CreatePurchaseOrderItemRequest struct {
	MedicineID uint    `json:"medicine_id" validate:"required"`
	Quantity   int     `json:"quantity" validate:"required,gt=0"`
	Price      float64 `json:"price" validate:"required,gt=0"`
}

type CreatePurchaseOrderRequest struct {
	SupplierID uint                             `json:"supplier_id" validate:"required"`
	Items      []CreatePurchaseOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}
