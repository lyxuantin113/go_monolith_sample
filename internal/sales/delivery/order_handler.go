package delivery

import (
	"errors"
	med "go_monolith_sample/internal/domain/medicine"
	"go_monolith_sample/internal/domain/sales"
	"go_monolith_sample/internal/sales/delivery/dto"
	apperror "go_monolith_sample/pkg/error"
	response "go_monolith_sample/pkg/response"
	validate "go_monolith_sample/pkg/validate"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	salesService sales.OrderService
}

func NewOrderHandler(service sales.OrderService) *OrderHandler {
	return &OrderHandler{salesService: service}
}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	var req dto.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, apperror.BadRequest("Dữ liệu không hợp lệ", err))
		return
	}

	if errs := validate.ValidateStruct(&req); errs != nil {
		response.ValidationErrors(ctx, errs)
		return
	}

	input := sales.CreateOrderInput{
		CustomerID: req.CustomerID,
		Items: make([]struct {
			MedicineID uint `json:"medicine_id" validate:"required"`
			Quantity   int  `json:"quantity" validate:"required,gt=0"`
		}, len(req.Items)),
	}

	for i, item := range req.Items {
		input.Items[i].MedicineID = item.MedicineID
		input.Items[i].Quantity = item.Quantity
	}

	// 4
	data, err := h.salesService.CreateOrder(ctx.Request.Context(), input)
	if err != nil {
		if errors.Is(err, med.ErrInsufficientStock) {
			response.Error(ctx, apperror.BadRequest("Hết hàng rồi cậu ơi", err))
			return
		}
		response.Error(ctx, apperror.Internal("Lỗi hệ thống khi tạo đơn", err))
		return
	}

	response.Success(ctx, "Đơn hàng đã được tạo thành công", data, nil)
}
