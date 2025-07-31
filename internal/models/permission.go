package models

import "time"

type Permission struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Codename    string    `json:"codename" gorm:"not null"`
	ContentType string    `json:"content_type" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AuthGroup struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name" gorm:"uniqueIndex;not null"`
	Permissions []Permission `json:"permissions" gorm:"many2many:auth_group_permissions;"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// AuthGroup DTOs and Requests

// AuthGroupRequest for creating/updating auth groups
type AuthGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Permissions []uint `json:"permissions"`
}

// PaginatedAuthGroupResponse for Swagger documentation
type PaginatedAuthGroupResponse struct {
	Data       []AuthGroup `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// PaginatedPermissionResponse for Swagger documentation
type PaginatedPermissionResponse struct {
	Data       []Permission `json:"data"`
	Total      int          `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
}
