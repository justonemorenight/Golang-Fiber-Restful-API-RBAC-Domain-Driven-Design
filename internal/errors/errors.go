package errors

import "github.com/gofiber/fiber/v2"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message string, detail string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// Common errors
var (
	ErrBadRequest        = NewAppError(fiber.StatusBadRequest, "Bad Request", "")
	ErrUnauthorized      = NewAppError(fiber.StatusUnauthorized, "Unauthorized", "")
	ErrForbidden         = NewAppError(fiber.StatusForbidden, "Forbidden", "")
	ErrNotFound          = NewAppError(fiber.StatusNotFound, "Not Found", "")
	ErrInternalServer    = NewAppError(fiber.StatusInternalServerError, "Internal Server Error", "")
	ErrValidation        = NewAppError(fiber.StatusBadRequest, "Validation Error", "")
	ErrDatabaseOperation = NewAppError(fiber.StatusInternalServerError, "Database Error", "")
)
