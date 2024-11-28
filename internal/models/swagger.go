package models

import "time"

// UserResponse là struct dùng cho Swagger documentation
type UserResponse struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginResponseSwagger là struct dùng cho Swagger documentation
type LoginResponseSwagger struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

// ErrorResponse là struct dùng cho Swagger documentation
type ErrorResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Detail  string `json:"detail,omitempty"`
	} `json:"error"`
}

// RefreshTokenRequest là struct dùng cho Swagger documentation
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse là struct dùng cho Swagger documentation
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
