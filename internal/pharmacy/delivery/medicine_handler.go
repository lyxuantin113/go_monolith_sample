package delivery

import (
	"go_monolith_sample/internal/domain"
	"go_monolith_sample/internal/pharmacy/delivery/dto"
	"go_monolith_sample/pkg/utils"
	"net/http"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errs := utils.ValidateStruct(medicine); errs != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	medicineDomain := &domain.Medicine{
		Name:        medicine.Name,
		Price:       medicine.Price,
		Stock:       medicine.Stock,
		Description: medicine.Description,
	}

	if err := h.medicineService.CreateMedicine(medicineDomain); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Thuốc đã được tạo thành công", "data": medicine})
}

func (h *MedicineHandler) UpdateMedicine(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdateMedicineRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errs := utils.ValidateStruct(req); errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	// Chuyển DTO thành Domain Input
	input := domain.UpdateMedicineInput{
		Name:        req.Name,
		Description: req.Description,
		Stock:       req.Stock,
		Price:       req.Price,
	}

	if err := h.medicineService.UpdateMedicine(uint(id), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

func (h *MedicineHandler) DeleteMedicine(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	if err := h.medicineService.DeleteMedicine(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Thuốc đã được xóa thành công"})
}

func (h *MedicineHandler) GetMedicineByID(ctx *gin.Context) {
	id := ctx.Param("id")

	uintID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	medicine, err := h.medicineService.GetMedicineByID(uint(uintID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Thuốc đã được lấy thành công", "data": medicine})
}

func (h *MedicineHandler) GetAllMedicines(ctx *gin.Context) {
	medicines, err := h.medicineService.GetAllMedicines()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Thuốc đã được lấy thành công", "data": medicines})
}
