package models

import (
	"time"

	"gorm.io/gorm"
)

// Organization model
type Organization struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null"`
	Domain    *string   `json:"domain,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Users     []User    `json:"users,omitempty" gorm:"foreignKey:OrganizationID"`
}

// Organization methods
func (o *Organization) BeforeCreate(tx *gorm.DB) error {
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Organization) BeforeUpdate(tx *gorm.DB) error {
	o.UpdatedAt = time.Now()
	return nil
}

// Organization DTOs and Requests

// OrganizationCreateRequest for creating organizations
type OrganizationCreateRequest struct {
	Name   string  `json:"name" binding:"required"`
	Domain *string `json:"domain,omitempty"`
}

// OrganizationUpdateRequest for updating organizations
type OrganizationUpdateRequest struct {
	Name   *string `json:"name,omitempty"`
	Domain *string `json:"domain,omitempty"`
}

// OrganizationResponse for API responses
type OrganizationResponse struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	Domain    *string `json:"domain,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// PaginatedOrganizationResponse for Swagger documentation
type PaginatedOrganizationResponse struct {
	Data       []OrganizationResponse `json:"data"`
	Total      int                    `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPages int                    `json:"total_pages"`
}
