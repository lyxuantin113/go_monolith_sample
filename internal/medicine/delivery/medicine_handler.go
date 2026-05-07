package delivery

import (
	med "go_monolith_sample/internal/domain/medicine"
	"go_monolith_sample/internal/medicine/delivery/dto"
	apperror "go_monolith_sample/pkg/error"
	response "go_monolith_sample/pkg/response"
	validate "go_monolith_sample/pkg/validate"
	"strconv"

	"github.com/gin-gonic/gin"
)

type medicineHandler struct {
	medicineService med.MedicineService
}

func NewMedicineHandler(medicineService med.MedicineService) *medicineHandler {
	return &medicineHandler{
		medicineService: medicineService,
	}
}

func (h *medicineHandler) CreateMedicine(ctx *gin.Context) {
	var medicine dto.CreateMedicineRequest
	if err := ctx.ShouldBindJSON(&medicine); err != nil {
		response.Error(ctx, apperror.BadRequest("Dữ liệu không hợp lệ", err))
		return
	}

	if errs := validate.ValidateStruct(medicine); errs != nil {
		response.ValidationErrors(ctx, errs)
		return
	}

	medicineDomain := &med.Medicine{
		Name:        medicine.Name,
		Price:       medicine.Price,
		Stock:       medicine.Stock,
		Description: medicine.Description,
	}

	if err := h.medicineService.CreateMedicine(ctx.Request.Context(), medicineDomain); err != nil {
		response.Error(ctx, apperror.Internal("Lỗi khi tạo thuốc", err))
		return
	}
	response.Success(ctx, "Thuốc đã được tạo thành công", medicine, nil)
}

func (h *medicineHandler) UpdateMedicine(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req dto.UpdateMedicineRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, apperror.BadRequest("Dữ liệu không hợp lệ", err))
		return
	}

	if errs := validate.ValidateStruct(req); errs != nil {
		response.ValidationErrors(ctx, errs)
		return
	}

	// Chuyển DTO thành Domain Input
	input := med.UpdateMedicineInput{
		Name:        req.Name,
		Description: req.Description,
		Stock:       req.Stock,
		Price:       req.Price,
	}

	if err := h.medicineService.UpdateMedicine(ctx.Request.Context(), uint(id), input); err != nil {
		response.Error(ctx, apperror.Internal("Lỗi khi cập nhật thuốc", err))
		return
	}

	response.Success(ctx, "Thuốc đã được cập nhật thành công", nil, nil)
}

func (h *medicineHandler) DeleteMedicine(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, apperror.BadRequest("ID không hợp lệ", err))
		return
	}

	if err := h.medicineService.DeleteMedicine(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, apperror.Internal("Lỗi khi xóa thuốc", err))
		return
	}

	response.Success(ctx, "Thuốc đã được xóa thành công", nil, nil)
}

func (h *medicineHandler) GetMedicineByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, apperror.BadRequest("ID không hợp lệ", err))
		return
	}

	medicine, err := h.medicineService.GetMedicineByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, apperror.Internal("Lỗi khi lấy thông tin thuốc", err))
		return
	}

	response.Success(ctx, "Thuốc đã được lấy thành công", medicine, nil)
}

func (h *medicineHandler) GetAllMedicines(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	search := ctx.Query("search")

	if err != nil {
		response.Error(ctx, apperror.BadRequest("Invalid page parameter", err))
		return
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if err != nil {
		response.Error(ctx, apperror.BadRequest("Invalid page_size parameter", err))
		return
	}

	medicines, pagination, err := h.medicineService.GetAllMedicines(ctx.Request.Context(), page, pageSize, search)
	if err != nil {
		response.Error(ctx, apperror.Internal("Lỗi khi lấy thông tin thuốc", err))
		return
	}
	response.Success(ctx, "Thuốc đã được lấy thành công", gin.H{"data": medicines, "pagination": pagination}, nil)
}
