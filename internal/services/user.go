package services

import (
	"errors"
	"kepler-auth-go/internal/database"
	"kepler-auth-go/internal/models"
	"math"

	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetUsers(query *models.PaginationQuery, organizationID *uint) (*models.PaginatedResponse[models.UserResponse], error) {
	var users []models.User
	var total int64

	db := database.GetDB().Model(&models.User{}).Preload("Groups").Preload("Organization")

	// Filter by organization
	if organizationID != nil {
		db = db.Where("organization_id = ?", *organizationID)
	}

	if query.Search != "" {
		db = db.Where("name ILIKE ? OR email ILIKE ? OR phone_number ILIKE ?",
			"%"+query.Search+"%", "%"+query.Search+"%", "%"+query.Search+"%")
	}

	if query.Status != "" {
		switch query.Status {
		case "active":
			db = db.Where("is_verified = ? AND is_active = ? AND is_deleted = ?", true, true, false)
		case "pending":
			db = db.Where("is_verified = ?", false)
		case "deactivated":
			db = db.Where("is_verified = ? AND (is_deleted = ? OR is_active = ?)", true, true, false)
		}
	}

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&users).Error; err != nil {
		return nil, err
	}

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = s.toUserResponse(&user)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &models.PaginatedResponse[models.UserResponse]{
		Data:       userResponses,
		Total:      int(total),
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) GetUserByID(id uint, organizationID *uint) (*models.UserResponse, error) {
	var user models.User
	db := database.GetDB().Preload("Groups").Preload("Organization")

	// Filter by organization if provided
	if organizationID != nil {
		db = db.Where("organization_id = ?", *organizationID)
	}

	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	response := s.toUserResponse(&user)
	return &response, nil
}

func (s *UserService) UpdateUser(id uint, req *models.UserUpdateRequest, organizationID *uint) (*models.UserResponse, error) {
	var user models.User
	db := database.GetDB()

	// Filter by organization if provided
	if organizationID != nil {
		db = db.Where("organization_id = ?", *organizationID)
	}

	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.PhoneNumber != nil {
		updates["phone_number"] = *req.PhoneNumber
	}
	if req.ProfilePicture != nil {
		updates["profile_picture"] = *req.ProfilePicture
	}
	if req.Country != nil {
		updates["country"] = *req.Country
	}
	if req.City != nil {
		updates["city"] = *req.City
	}
	if req.WhatsappNo != nil {
		updates["whatsapp_no"] = *req.WhatsappNo
	}
	if req.SendWhatsapp != nil {
		updates["send_whatsapp"] = *req.SendWhatsapp
	}
	if req.SendEmail != nil {
		updates["send_email"] = *req.SendEmail
	}
	if req.OrganizationID != nil {
		// Validate organization exists
		var org models.Organization
		if err := database.GetDB().First(&org, *req.OrganizationID).Error; err != nil {
			return nil, errors.New("organization not found")
		}
		updates["organization_id"] = *req.OrganizationID
	}

	if err := database.GetDB().Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := database.GetDB().Preload("Groups").Preload("Organization").First(&user, id).Error; err != nil {
		return nil, err
	}

	response := s.toUserResponse(&user)
	return &response, nil
}

func (s *UserService) DeleteUser(id uint, organizationID *uint) error {
	var user models.User
	db := database.GetDB()

	// Filter by organization if provided
	if organizationID != nil {
		db = db.Where("organization_id = ?", *organizationID)
	}

	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return database.GetDB().Model(&user).Update("is_deleted", true).Error
}

func (s *UserService) toUserResponse(user *models.User) models.UserResponse {
	// Collect all permissions from user's groups
	permissions := make([]int, 0)
	for _, group := range user.Groups {
		permissions = append(permissions, group.Permissions...)
	}

	// Remove duplicates
	permissionSet := make(map[int]bool)
	uniquePermissions := make([]int, 0)
	for _, perm := range permissions {
		if !permissionSet[perm] {
			permissionSet[perm] = true
			uniquePermissions = append(uniquePermissions, perm)
		}
	}

	return models.UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		Name:           user.Name,
		PhoneNumber:    user.PhoneNumber,
		ProfilePicture: user.ProfilePicture,
		Country:        user.Country,
		City:           user.City,
		WhatsappNo:     user.WhatsappNo,
		SendWhatsapp:   user.SendWhatsapp,
		SendEmail:      user.SendEmail,
		IsVerified:     user.IsVerified,
		IsDeleted:      user.IsDeleted,
		IsStaff:        user.IsStaff,
		IsAdmin:        user.IsAdmin,
		IsActive:       user.IsActive,
		OrganizationID: user.OrganizationID,
		Organization:   user.Organization,
		Status:         user.GetStatus(),
		Groups:         user.Groups,
		Permissions:    uniquePermissions,
	}
}
