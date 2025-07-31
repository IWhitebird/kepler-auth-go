package services

import (
	"kepler-auth-go/internal/database"
	"kepler-auth-go/internal/models"
	"math"
)

type PermissionService struct{}

func NewPermissionService() *PermissionService {
	return &PermissionService{}
}

func (s *PermissionService) GetPermissions(query *models.PaginationQuery) (*models.PaginatedPermissionResponse, error) {
	var permissions []models.Permission
	var total int64

	db := database.GetDB().Model(&models.Permission{})

	if query.Search != "" {
		db = db.Where("name ILIKE ? OR codename ILIKE ? OR content_type ILIKE ?",
			"%"+query.Search+"%", "%"+query.Search+"%", "%"+query.Search+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&permissions).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &models.PaginatedPermissionResponse{
		Data:       permissions,
		Total:      int(total),
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *PermissionService) GetAuthGroups(query *models.PaginationQuery) (*models.PaginatedAuthGroupResponse, error) {
	var authGroups []models.AuthGroup
	var total int64

	db := database.GetDB().Model(&models.AuthGroup{}).Preload("Permissions")

	if query.Search != "" {
		db = db.Where("name ILIKE ?", "%"+query.Search+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&authGroups).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &models.PaginatedAuthGroupResponse{
		Data:       authGroups,
		Total:      int(total),
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *PermissionService) CreateAuthGroup(req *models.AuthGroupRequest) (*models.AuthGroup, error) {
	authGroup := &models.AuthGroup{
		Name: req.Name,
	}

	if err := database.GetDB().Create(authGroup).Error; err != nil {
		return nil, err
	}

	if len(req.Permissions) > 0 {
		var permissions []models.Permission
		if err := database.GetDB().Where("id IN ?", req.Permissions).Find(&permissions).Error; err != nil {
			return nil, err
		}
		if err := database.GetDB().Model(authGroup).Association("Permissions").Replace(&permissions); err != nil {
			return nil, err
		}
	}

	if err := database.GetDB().Preload("Permissions").First(authGroup, authGroup.ID).Error; err != nil {
		return nil, err
	}

	return authGroup, nil
}

func (s *PermissionService) UpdateAuthGroup(id uint, req *models.AuthGroupRequest) (*models.AuthGroup, error) {
	var authGroup models.AuthGroup
	if err := database.GetDB().First(&authGroup, id).Error; err != nil {
		return nil, err
	}

	authGroup.Name = req.Name
	if err := database.GetDB().Save(&authGroup).Error; err != nil {
		return nil, err
	}

	var permissions []models.Permission
	if len(req.Permissions) > 0 {
		if err := database.GetDB().Where("id IN ?", req.Permissions).Find(&permissions).Error; err != nil {
			return nil, err
		}
	}

	if err := database.GetDB().Model(&authGroup).Association("Permissions").Replace(&permissions); err != nil {
		return nil, err
	}

	if err := database.GetDB().Preload("Permissions").First(&authGroup, id).Error; err != nil {
		return nil, err
	}

	return &authGroup, nil
}

func (s *PermissionService) DeleteAuthGroup(id uint) error {
	return database.GetDB().Delete(&models.AuthGroup{}, id).Error
}
