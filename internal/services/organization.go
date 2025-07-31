package services

import (
	"errors"
	"kepler-auth-go/internal/database"
	"kepler-auth-go/internal/models"
	"math"
	"time"

	"gorm.io/gorm"
)

type OrganizationService struct{}

func NewOrganizationService() *OrganizationService {
	return &OrganizationService{}
}

func (s *OrganizationService) GetOrganizations(query *models.PaginationQuery) (*models.PaginatedResponse[models.OrganizationResponse], error) {
	var organizations []models.Organization
	var total int64

	db := database.GetDB().Model(&models.Organization{})

	if query.Search != "" {
		db = db.Where("name ILIKE ? OR domain ILIKE ?",
			"%"+query.Search+"%", "%"+query.Search+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&organizations).Error; err != nil {
		return nil, err
	}

	organizationResponses := make([]models.OrganizationResponse, len(organizations))
	for i, org := range organizations {
		organizationResponses[i] = s.toOrganizationResponse(&org)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &models.PaginatedResponse[models.OrganizationResponse]{
		Data:       organizationResponses,
		Total:      int(total),
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *OrganizationService) GetOrganizationByID(id uint) (*models.OrganizationResponse, error) {
	var organization models.Organization
	if err := database.GetDB().First(&organization, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organization not found")
		}
		return nil, err
	}

	response := s.toOrganizationResponse(&organization)
	return &response, nil
}

func (s *OrganizationService) CreateOrganization(req *models.OrganizationCreateRequest) (*models.OrganizationResponse, error) {
	// Check if organization with same name exists
	var existingOrg models.Organization
	if err := database.GetDB().Where("name = ?", req.Name).First(&existingOrg).Error; err == nil {
		return nil, errors.New("organization with this name already exists")
	}

	organization := &models.Organization{
		Name:   req.Name,
		Domain: req.Domain,
	}

	if err := database.GetDB().Create(organization).Error; err != nil {
		return nil, err
	}

	response := s.toOrganizationResponse(organization)
	return &response, nil
}

func (s *OrganizationService) UpdateOrganization(id uint, req *models.OrganizationUpdateRequest) (*models.OrganizationResponse, error) {
	var organization models.Organization
	if err := database.GetDB().First(&organization, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organization not found")
		}
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		// Check if another organization with same name exists
		var existingOrg models.Organization
		if err := database.GetDB().Where("name = ? AND id != ?", *req.Name, id).First(&existingOrg).Error; err == nil {
			return nil, errors.New("organization with this name already exists")
		}
		updates["name"] = *req.Name
	}
	if req.Domain != nil {
		updates["domain"] = *req.Domain
	}

	if err := database.GetDB().Model(&organization).Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := database.GetDB().First(&organization, id).Error; err != nil {
		return nil, err
	}

	response := s.toOrganizationResponse(&organization)
	return &response, nil
}

func (s *OrganizationService) DeleteOrganization(id uint) error {
	var organization models.Organization
	if err := database.GetDB().First(&organization, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("organization not found")
		}
		return err
	}

	// Check if organization has users
	var userCount int64
	if err := database.GetDB().Model(&models.User{}).Where("organization_id = ?", id).Count(&userCount).Error; err != nil {
		return err
	}

	if userCount > 0 {
		return errors.New("cannot delete organization with existing users")
	}

	return database.GetDB().Delete(&organization).Error
}

func (s *OrganizationService) toOrganizationResponse(org *models.Organization) models.OrganizationResponse {
	return models.OrganizationResponse{
		ID:        org.ID,
		Name:      org.Name,
		Domain:    org.Domain,
		CreatedAt: org.CreatedAt.Format(time.RFC3339),
		UpdatedAt: org.UpdatedAt.Format(time.RFC3339),
	}
}
