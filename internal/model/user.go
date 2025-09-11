package model

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"unique;not null;size:50" json:"name" binding:"required,oneof=user admin"`
	Description string         `gorm:"size:255" json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"unique;not null" json:"username" binding:"required,min=3"`
	Email     string         `gorm:"unique;not null" json:"email" binding:"required,email"`
	Password  string         `gorm:"not null" json:"-"` // Exclude from JSON
	RoleID    uint           `gorm:"not null" json:"role_id" binding:"required"`
	Role      Role           `gorm:"foreignKey:RoleID" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}
