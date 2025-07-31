package models

import (
	"time"

	"gorm.io/gorm"
)

// User model
type User struct {
	ID             uint          `json:"id" gorm:"primaryKey"`
	Email          string        `json:"email" gorm:"not null;index:idx_email_org,unique"`
	Name           string        `json:"name" gorm:"not null"`
	Password       string        `json:"-" gorm:"not null"`
	PhoneNumber    *string       `json:"phone_number,omitempty"`
	ProfilePicture *string       `json:"profile_picture,omitempty"`
	Country        *string       `json:"country,omitempty"`
	City           *string       `json:"city,omitempty"`
	WhatsappNo     *string       `json:"whatsapp_no,omitempty"`
	SendWhatsapp   bool          `json:"send_whatsapp" gorm:"default:false"`
	SendEmail      bool          `json:"send_email" gorm:"default:false"`
	IsVerified     bool          `json:"is_verified" gorm:"default:false"`
	IsDeleted      bool          `json:"is_deleted" gorm:"default:false"`
	IsStaff        bool          `json:"is_staff" gorm:"default:false"`
	IsAdmin        bool          `json:"is_admin" gorm:"default:false"`
	IsActive       bool          `json:"is_active" gorm:"default:true"`
	OrganizationID *uint         `json:"organization_id,omitempty" gorm:"index:idx_email_org,unique;index"`
	Organization   *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	Groups         []Group       `json:"groups,omitempty" gorm:"many2many:user_groups;"`
}

// Group model
type Group struct {
	ID             uint          `json:"id" gorm:"primaryKey"`
	Name           string        `json:"name" gorm:"not null;index:idx_group_name_org,unique"`
	Description    *string       `json:"description,omitempty"`
	Permissions    []int         `json:"permissions" gorm:"type:integer[]"`
	IsActive       bool          `json:"is_active" gorm:"default:true"`
	IsDefault      bool          `json:"is_default" gorm:"default:false"`
	OrganizationID *uint         `json:"organization_id,omitempty" gorm:"index:idx_group_name_org,unique;index"`
	Organization   *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	Users          []User        `json:"users,omitempty" gorm:"many2many:user_groups;"`
}

// UserStatus enum
type UserStatus string

const (
	StatusPending     UserStatus = "pending"
	StatusActive      UserStatus = "active"
	StatusDeactivated UserStatus = "deactivated"
	StatusUnknown     UserStatus = "unknown"
)

// User methods
func (u *User) GetStatus() UserStatus {
	if !u.IsVerified {
		return StatusPending
	}
	if u.IsVerified && (u.IsDeleted || !u.IsActive) {
		return StatusDeactivated
	}
	if u.IsVerified && u.IsActive && !u.IsDeleted {
		return StatusActive
	}
	return StatusUnknown
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// User DTOs and Requests

// RegisterRequest for user registration
type RegisterRequest struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	OrganizationID *uint  `json:"organization_id,omitempty"`
}

// LoginRequest for user authentication
type LoginRequest struct {
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required"`
	OrganizationID *uint  `json:"organization_id,omitempty"`
}

// LoginResponse after successful authentication
type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// ChangePasswordRequest for password updates
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ResetPasswordEmailRequest for password reset emails
type ResetPasswordEmailRequest struct {
	Email       string `json:"email" binding:"required,email"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

// SetNewPasswordRequest for setting new password after reset
type SetNewPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	UID      string `json:"uid" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserUpdateRequest for user profile updates
type UserUpdateRequest struct {
	Name           *string `json:"name,omitempty"`
	PhoneNumber    *string `json:"phone_number,omitempty"`
	ProfilePicture *string `json:"profile_picture,omitempty"`
	Country        *string `json:"country,omitempty"`
	City           *string `json:"city,omitempty"`
	WhatsappNo     *string `json:"whatsapp_no,omitempty"`
	SendWhatsapp   *bool   `json:"send_whatsapp,omitempty"`
	SendEmail      *bool   `json:"send_email,omitempty"`
	OrganizationID *uint   `json:"organization_id,omitempty"`
}

// UserResponse for API responses
type UserResponse struct {
	ID             uint          `json:"id"`
	Email          string        `json:"email"`
	Name           string        `json:"name"`
	PhoneNumber    *string       `json:"phone_number,omitempty"`
	ProfilePicture *string       `json:"profile_picture,omitempty"`
	Country        *string       `json:"country,omitempty"`
	City           *string       `json:"city,omitempty"`
	WhatsappNo     *string       `json:"whatsapp_no,omitempty"`
	SendWhatsapp   bool          `json:"send_whatsapp"`
	SendEmail      bool          `json:"send_email"`
	IsVerified     bool          `json:"is_verified"`
	IsDeleted      bool          `json:"is_deleted"`
	IsStaff        bool          `json:"is_staff"`
	IsAdmin        bool          `json:"is_admin"`
	IsActive       bool          `json:"is_active"`
	OrganizationID *uint         `json:"organization_id,omitempty"`
	Organization   *Organization `json:"organization,omitempty"`
	Status         UserStatus    `json:"status"`
	Groups         []Group       `json:"groups,omitempty"`
	Permissions    []int         `json:"permissions,omitempty"`
}

// Group DTOs and Requests

// GroupRequest for creating/updating groups
type GroupRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	Permissions []int   `json:"permissions"`
	IsActive    *bool   `json:"is_active,omitempty"`
	IsDefault   *bool   `json:"is_default,omitempty"`
}

// PaginatedUserResponse for Swagger documentation
type PaginatedUserResponse struct {
	Data       []UserResponse `json:"data"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// PaginatedGroupResponse for Swagger documentation
type PaginatedGroupResponse struct {
	Data       []Group `json:"data"`
	Total      int     `json:"total"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	TotalPages int     `json:"total_pages"`
}
