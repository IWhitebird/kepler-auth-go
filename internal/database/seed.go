package database

import (
	"kepler-auth-go/internal/models"
	"log"
)

func SeedDefaultData() error {
	// Seed default permissions (similar to Django's default permissions)
	permissions := []models.Permission{
		{Name: "Can add user", Codename: "add_user", ContentType: "auth.user"},
		{Name: "Can change user", Codename: "change_user", ContentType: "auth.user"},
		{Name: "Can delete user", Codename: "delete_user", ContentType: "auth.user"},
		{Name: "Can view user", Codename: "view_user", ContentType: "auth.user"},
		{Name: "Can add group", Codename: "add_group", ContentType: "auth.group"},
		{Name: "Can change group", Codename: "change_group", ContentType: "auth.group"},
		{Name: "Can delete group", Codename: "delete_group", ContentType: "auth.group"},
		{Name: "Can view group", Codename: "view_group", ContentType: "auth.group"},
		{Name: "Can add permission", Codename: "add_permission", ContentType: "auth.permission"},
		{Name: "Can change permission", Codename: "change_permission", ContentType: "auth.permission"},
		{Name: "Can delete permission", Codename: "delete_permission", ContentType: "auth.permission"},
		{Name: "Can view permission", Codename: "view_permission", ContentType: "auth.permission"},
	}

	for _, permission := range permissions {
		var existing models.Permission
		if err := DB.Where("codename = ? AND content_type = ?", permission.Codename, permission.ContentType).First(&existing).Error; err != nil {
			if err := DB.Create(&permission).Error; err != nil {
				log.Printf("Failed to create permission %s: %v", permission.Name, err)
			} else {
				log.Printf("Created permission: %s", permission.Name)
			}
		}
	}

	// Seed default groups
	defaultGroups := []models.AuthGroup{
		{Name: "Administrators"},
		{Name: "Staff"},
		{Name: "Users"},
	}

	for _, group := range defaultGroups {
		var existing models.AuthGroup
		if err := DB.Where("name = ?", group.Name).First(&existing).Error; err != nil {
			if err := DB.Create(&group).Error; err != nil {
				log.Printf("Failed to create auth group %s: %v", group.Name, err)
			} else {
				log.Printf("Created auth group: %s", group.Name)
			}
		}
	}

	// Seed default custom groups
	customGroups := []models.Group{
		{Name: "Admin", Description: stringPtr("Administrator group with full access"), IsActive: true, IsDefault: false},
		{Name: "Staff", Description: stringPtr("Staff group with limited access"), IsActive: true, IsDefault: false},
		{Name: "User", Description: stringPtr("Default user group"), IsActive: true, IsDefault: true},
	}

	for _, group := range customGroups {
		var existing models.Group
		if err := DB.Where("name = ?", group.Name).First(&existing).Error; err != nil {
			if err := DB.Create(&group).Error; err != nil {
				log.Printf("Failed to create group %s: %v", group.Name, err)
			} else {
				log.Printf("Created group: %s", group.Name)
			}
		}
	}

	return nil
}

func stringPtr(s string) *string {
	return &s
}
