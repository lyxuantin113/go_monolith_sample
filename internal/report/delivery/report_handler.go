package delivery

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go_monolith_sample/internal/domain/report"
	apperror "go_monolith_sample/pkg/error"
	"go_monolith_sample/pkg/response"

	"github.com/gin-gonic/gin"
)

type reportHandler struct {
	reportService report.ReportService
}

func NewReportHandler(service report.ReportService) *reportHandler {
	return &reportHandler{reportService: service}
}

// 1. API lấy Thẻ Kho (Dữ liệu JSON)
func (h *reportHandler) GetInventoryLedger(c *gin.Context) {
	medicineID, _ := strconv.ParseUint(c.Query("medicine_id"), 10, 32)
	from, _ := time.Parse("2006-01-02", c.Query("from"))
	to, _ := time.Parse("2006-01-02", c.Query("to"))

	if medicineID == 0 {
		response.Error(c, apperror.BadRequest("Thiếu Medicine ID", nil))
		return
	}

	data, err := h.reportService.GetInventoryLedger(c.Request.Context(), uint(medicineID), from, to)
	if err != nil {
		response.Error(c, apperror.Internal("Lỗi lấy báo cáo thẻ kho", err))
		return
	}

	response.Success(c, "Lấy thẻ kho thành công", data, nil)
}

// 2. API lấy Xuất Nhập Tồn (Dữ liệu JSON)
func (h *reportHandler) GetSIOReport(c *gin.Context) {
	from, _ := time.Parse("2006-01-02", c.Query("from"))
	to, _ := time.Parse("2006-01-02", c.Query("to"))

	data, err := h.reportService.GetSIOReport(c.Request.Context(), from, to)
	if err != nil {
		response.Error(c, apperror.Internal("Lỗi lấy báo cáo SIO", err))
		return
	}

	response.Success(c, "Lấy báo cáo Xuất-Nhập-Tồn thành công", data, nil)
}

// 3. API Xuất Excel
func (h *reportHandler) ExportSIOExcel(c *gin.Context) {
	from, _ := time.Parse("2006-01-02", c.Query("from"))
	to, _ := time.Parse("2006-01-02", c.Query("to"))

	fileBytes, err := h.reportService.ExportSIOToExcel(c.Request.Context(), from, to)
	if err != nil {
		response.Error(c, apperror.Internal("Lỗi xuất file excel", err))
		return
	}

	fileName := fmt.Sprintf("Bao-Cao-SIO-%s.xlsx", time.Now().Format("2006-01-02"))

	// Thiết lập các Header quan trọng để Browser kích hoạt chế độ tải file
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")

	c.Data(http.StatusOK, "application/octet-stream", fileBytes)
}

// 4. API Chụp ảnh kho thủ công (Dành cho Admin)
func (h *reportHandler) TakeSnapshot(c *gin.Context) {
	err := h.reportService.TakeInventorySnapshot(c.Request.Context())
	if err != nil {
		response.Error(c, apperror.Internal("Lỗi chụp ảnh kho", err))
		return
	}

	response.Success(c, "Chụp ảnh kho thành công", nil, nil)
}
