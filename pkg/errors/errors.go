package errors

import "net/http"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrInvalidExpression = NewAppError(http.StatusUnprocessableEntity, "Expression is not valid")
	ErrInternalServer    = NewAppError(http.StatusInternalServerError, "Internal server error")
	ErrBadRequest        = NewAppError(http.StatusBadRequest, "Invalid JSON")
)
