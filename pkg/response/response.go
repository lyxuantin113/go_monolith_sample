package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperror "go_monolith_sample/pkg/error"
)

type JSONResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func Success(ctx *gin.Context, message string, data interface{}, meta interface{}) {
	ctx.JSON(http.StatusOK, JSONResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(ctx *gin.Context, err error) {
	if appErr, ok := err.(*apperror.AppError); ok {
		ctx.JSON(appErr.StatusCode, JSONResponse{
			Success: false,
			Message: appErr.Message,
		})
		return
	}

	ctx.JSON(http.StatusInternalServerError, JSONResponse{
		Success: false,
		Message: "Internal Server Error",
	})
}

func ValidationErrors(ctx *gin.Context, errs interface{}) {
	ctx.JSON(http.StatusBadRequest, JSONResponse{
		Success: false,
		Message: "Dữ liệu không hợp lệ",
		Data:    errs,
	})
}
