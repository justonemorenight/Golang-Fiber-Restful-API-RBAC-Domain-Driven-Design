package models

import "time"

// User related
type SwaggerUserResponse struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SwaggerLoginResponse struct {
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	User         SwaggerUserResponse `json:"user"`
}

// Role related
type SwaggerRoleResponse struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SwaggerCreateRoleRequest struct {
	Name        string `json:"name" validate:"required" example:"admin"`
	Description string `json:"description" example:"Administrator role"`
}

type SwaggerUpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Common response wrapper for all API endpoints
type SwaggerResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorData  `json:"error,omitempty"`
}

// Error data structure
type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// ErrorResponse represents the error response structure
type ErrorResponse struct {
	Success bool      `json:"success"`
	Error   ErrorData `json:"error"`
}

// Common success response helper
func NewSuccessResponse(data interface{}) *SwaggerResponse {
	return &SwaggerResponse{
		Success: true,
		Data:    data,
	}
}

// Common error response helper
func NewErrorResponse(code int, message string, detail string) *SwaggerResponse {
	return &SwaggerResponse{
		Success: false,
		Error: &ErrorData{
			Code:    code,
			Message: message,
			Detail:  detail,
		},
	}
}
