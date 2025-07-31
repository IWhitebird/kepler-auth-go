package services

import (
	"errors"
	"kepler-auth-go/internal/database"
	"kepler-auth-go/internal/models"
	"math"

	"gorm.io/gorm"
)

type GroupService struct{}

func NewGroupService() *GroupService {
	return &GroupService{}
}

func (s *GroupService) GetGroups(query *models.PaginationQuery) (*models.PaginatedGroupResponse, error) {
	var groups []models.Group
	var total int64

	db := database.GetDB().Model(&models.Group{})

	if query.Search != "" {
		db = db.Where("name ILIKE ? OR description ILIKE ?", "%"+query.Search+"%", "%"+query.Search+"%")
	}

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&groups).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &models.PaginatedGroupResponse{
		Data:       groups,
		Total:      int(total),
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *GroupService) GetGroupByID(id uint) (*models.Group, error) {
	var group models.Group
	if err := database.GetDB().First(&group, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("group not found")
		}
		return nil, err
	}
	return &group, nil
}

func (s *GroupService) CreateGroup(req *models.GroupRequest) (*models.Group, error) {
	var existingGroup models.Group
	if err := database.GetDB().Where("name = ?", req.Name).First(&existingGroup).Error; err == nil {
		return nil, errors.New("group with this name already exists")
	}

	group := &models.Group{
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
		IsActive:    true,
		IsDefault:   false,
	}

	if req.IsActive != nil {
		group.IsActive = *req.IsActive
	}
	if req.IsDefault != nil {
		group.IsDefault = *req.IsDefault
	}

	if err := database.GetDB().Create(group).Error; err != nil {
		return nil, err
	}

	return group, nil
}

func (s *GroupService) UpdateGroup(id uint, req *models.GroupRequest) (*models.Group, error) {
	var group models.Group
	if err := database.GetDB().First(&group, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("group not found")
		}
		return nil, err
	}

	if req.Name != group.Name {
		var existingGroup models.Group
		if err := database.GetDB().Where("name = ? AND id != ?", req.Name, id).First(&existingGroup).Error; err == nil {
			return nil, errors.New("group with this name already exists")
		}
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"permissions": req.Permissions,
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.IsDefault != nil {
		updates["is_default"] = *req.IsDefault
	}

	if err := database.GetDB().Model(&group).Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := database.GetDB().First(&group, id).Error; err != nil {
		return nil, err
	}

	return &group, nil
}

func (s *GroupService) DeleteGroup(id uint) error {
	var group models.Group
	if err := database.GetDB().First(&group, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("group not found")
		}
		return err
	}

	return database.GetDB().Delete(&group).Error
}
