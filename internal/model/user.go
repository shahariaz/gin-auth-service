package model

import (
	"time"

	"gorm.io/gorm"
)

// Role represents a user role in the system
// @Description User role information
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id" example:"1"`
	Name        string         `gorm:"unique;not null;size:50" json:"name" binding:"required,oneof=user admin" example:"user"`
	Description string         `gorm:"size:255" json:"description,omitempty" example:"Regular user role"`
	CreatedAt   time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// User represents a user in the system
// @Description User account information
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id" example:"1"`
	Username  string         `gorm:"unique;not null" json:"username" binding:"required,min=3" example:"john_doe"`
	Email     string         `gorm:"unique;not null" json:"email" binding:"required,email" example:"john@example.com"`
	Password  string         `gorm:"not null" json:"-"` // Exclude from JSON
	RoleID    uint           `gorm:"not null" json:"role_id" binding:"required" example:"1"`
	Role      Role           `gorm:"foreignKey:RoleID" json:"role"`
	CreatedAt time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

// LoginRequest represents the login request payload
// @Description Login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// RegisterRequest represents the registration request payload
// @Description Registration request payload
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3" example:"john_doe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"password123"`
}

// RefreshTokenRequest represents the refresh token request payload
// @Description Refresh token request payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// UpdateProfileRequest represents the profile update request payload
// @Description Profile update request payload
type UpdateProfileRequest struct {
	Email string `json:"email" binding:"required,email" example:"newemail@example.com"`
}

// CreateUserRequest represents the admin user creation request payload
// @Description Admin user creation request payload
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3" example:"jane_doe"`
	Email    string `json:"email" binding:"required,email" example:"jane@example.com"`
	RoleID   uint   `json:"role_id" binding:"required" example:"1"`
}

// UpdateUserRequest represents the admin user update request payload
// @Description Admin user update request payload
type UpdateUserRequest struct {
	Username string `json:"username" binding:"required,min=3" example:"jane_doe_updated"`
	Email    string `json:"email" binding:"required,email" example:"jane.updated@example.com"`
	RoleID   uint   `json:"role_id" binding:"required" example:"2"`
}
