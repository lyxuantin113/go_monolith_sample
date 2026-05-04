package error

import "net/http"

type AppError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	RawError   error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

// Các hàm tạo lỗi nhanh
func BadRequest(msg string, err error) *AppError {
	return &AppError{StatusCode: http.StatusBadRequest, Message: msg, RawError: err}
}

func NotFound(msg string, err error) *AppError {
	return &AppError{StatusCode: http.StatusNotFound, Message: msg, RawError: err}
}

func Internal(msg string, err error) *AppError {
	return &AppError{StatusCode: http.StatusInternalServerError, Message: msg, RawError: err}
}
