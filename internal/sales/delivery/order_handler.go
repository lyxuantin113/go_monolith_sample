package delivery

import (
	"errors"
	med "go_monolith_sample/internal/domain/medicine"
	"go_monolith_sample/internal/domain/sales"
	"go_monolith_sample/internal/sales/delivery/dto"
	apperror "go_monolith_sample/pkg/error"
	response "go_monolith_sample/pkg/response"
	validate "go_monolith_sample/pkg/validate"
	"strconv"

	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	salesService sales.OrderService
}

func NewOrderHandler(service sales.OrderService) *orderHandler {
	return &orderHandler{salesService: service}
}

func (h *orderHandler) CreateOrder(ctx *gin.Context) {
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

func (h *orderHandler) RefundOrder(ctx *gin.Context) {
	// 1. Lấy ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, apperror.BadRequest("ID đơn hàng không hợp lệ", err))
		return
	}

	// 2. Gọi Service xử lý hoàn tiền
	err = h.salesService.RefundOrder(ctx.Request.Context(), uint(id))
	if err != nil {
		// Cậu có thể tùy biến thêm các loại lỗi cụ thể ở đây
		response.Error(ctx, apperror.Internal("Không thể hoàn tiền cho đơn hàng này", err))
		return
	}

	// 3. Trả về thành công
	response.Success(ctx, "Đã hoàn tiền và cập nhật kho thành công", nil, nil)
}
