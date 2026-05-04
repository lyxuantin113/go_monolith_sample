package delivery

import (
	"go_monolith_sample/internal/domain"
	"go_monolith_sample/internal/pharmacy/delivery/dto"
	apperror "go_monolith_sample/pkg/error"
	response "go_monolith_sample/pkg/response"
	utils "go_monolith_sample/pkg/validate"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MedicineHandler struct {
	medicineService domain.MedicineService
}

func NewMedicineHandler(medicineService domain.MedicineService) *MedicineHandler {
	return &MedicineHandler{
		medicineService: medicineService,
	}
}

func (h *MedicineHandler) CreateMedicine(ctx *gin.Context) {
	var medicine dto.CreateMedicineRequest
	if err := ctx.ShouldBindJSON(&medicine); err != nil {
		response.Error(ctx, apperror.BadRequest("Dữ liệu không hợp lệ", err))
		return
	}

	if errs := utils.ValidateStruct(medicine); errs != nil {
		response.ValidationErrors(ctx, errs)
		return
	}

	medicineDomain := &domain.Medicine{
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

func (h *MedicineHandler) UpdateMedicine(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req dto.UpdateMedicineRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, apperror.BadRequest("Dữ liệu không hợp lệ", err))
		return
	}

	if errs := utils.ValidateStruct(req); errs != nil {
		response.ValidationErrors(ctx, errs)
		return
	}

	// Chuyển DTO thành Domain Input
	input := domain.UpdateMedicineInput{
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

func (h *MedicineHandler) DeleteMedicine(ctx *gin.Context) {
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

func (h *MedicineHandler) GetMedicineByID(ctx *gin.Context) {
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

func (h *MedicineHandler) GetAllMedicines(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		response.Error(ctx, apperror.BadRequest("Invalid page parameter", err))
		return
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if err != nil {
		response.Error(ctx, apperror.BadRequest("Invalid page_size parameter", err))
		return
	}

	medicines, pagination, err := h.medicineService.GetAllMedicines(ctx.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(ctx, apperror.Internal("Lỗi khi lấy thông tin thuốc", err))
		return
	}
	response.Success(ctx, "Thuốc đã được lấy thành công", gin.H{"data": medicines, "pagination": pagination}, nil)
}
