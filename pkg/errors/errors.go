package errors

import "net/http"

// AppError - структура для представления ошибок приложения.
type AppError struct {
	Code    int    `json:"code"`    // HTTP код ошибки
	Message string `json:"message"` // Сообщение об ошибке
}

// NewAppError - создает новую ошибку приложения.
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Error - реализация метода Error() для интерфейса error.
func (e *AppError) Error() string {
	return e.Message
}

// Predefined errors
var (
	ErrInvalidExpression = NewAppError(http.StatusUnprocessableEntity, "Expression is not valid")
	ErrInternalServer    = NewAppError(http.StatusInternalServerError, "Internal server error")
	ErrBadRequest        = NewAppError(http.StatusBadRequest, "Invalid JSON")
)
