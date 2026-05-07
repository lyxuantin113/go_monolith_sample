package delivery

import (
	"errors"
	med "go_monolith_sample/internal/domain/medicine"
	"go_monolith_sample/internal/domain/purchase"
	"go_monolith_sample/internal/purchase/delivery/dto"
	apperror "go_monolith_sample/pkg/error"
	response "go_monolith_sample/pkg/response"
	validate "go_monolith_sample/pkg/validate"
	"strconv"

	"github.com/gin-gonic/gin"
)

type purchaseOrderHandler struct {
	purchaseOrderService purchase.PurchaseOrderService
}

func NewPurchaseOrderHandler(purchaseOrderService purchase.PurchaseOrderService) *purchaseOrderHandler {
	return &purchaseOrderHandler{
		purchaseOrderService: purchaseOrderService,
	}
}

func (h *purchaseOrderHandler) CreatePurchaseOrder(ctx *gin.Context) {

	var req dto.CreatePurchaseOrderRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, apperror.BadRequest("Dữ liệu không hợp lệ", err))
		return
	}

	if errs := validate.ValidateStruct(&req); errs != nil {
		response.ValidationErrors(ctx, errs)
		return
	}

	input := purchase.CreatePurchaseOrderInput{
		SupplierID: req.SupplierID,
		Items: make([]struct {
			MedicineID uint    `json:"medicine_id" binding:"required" validate:"required"`
			Quantity   int     `json:"quantity" binding:"required" validate:"required,gt=0"`
			Price      float64 `json:"price" binding:"required" validate:"required,gt=0"`
		}, len(req.Items)),
	}

	for i, item := range req.Items {
		input.Items[i].MedicineID = item.MedicineID
		input.Items[i].Quantity = item.Quantity
		input.Items[i].Price = item.Price
	}

	data, err := h.purchaseOrderService.CreatePurchaseOrder(ctx.Request.Context(), input)
	if err != nil {
		if errors.Is(err, med.ErrInsufficientStock) {
			response.Error(ctx, apperror.BadRequest("Hết hàng rồi cậu ơi", err))
			return
		}
		response.Error(ctx, apperror.Internal("Lỗi hệ thống khi tạo đơn", err))
		return
	}

	response.Success(ctx, "Đơn đặt hàng đã được tạo thành công", data, nil)
}

func (h *purchaseOrderHandler) CompletePurchaseOrder(ctx *gin.Context) {
	// Lấy ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, apperror.BadRequest("ID đơn nhập hàng không hợp lệ", err))
		return
	}

	// Gọi Service
	err = h.purchaseOrderService.CompletePurchaseOrder(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, apperror.Internal("Không thể hoàn thành đơn nhập này", err))
		return
	}

	response.Success(ctx, "Đã nhập hàng vào kho thành công", nil, nil)
}

func (h *purchaseOrderHandler) CancelPurchaseOrder(ctx *gin.Context) {
	// Lấy ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, apperror.BadRequest("ID đơn nhập hàng không hợp lệ", err))
		return
	}

	// Gọi Service
	err = h.purchaseOrderService.CancelPurchaseOrder(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, apperror.Internal("Không thể hủy đơn nhập này", err))
		return
	}

	response.Success(ctx, "Đã hủy đơn nhập hàng", nil, nil)
}
