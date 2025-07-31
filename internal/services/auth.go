package services

import (
	"errors"
	"kepler-auth-go/internal/config"
	"kepler-auth-go/internal/database"
	"kepler-auth-go/internal/middleware"
	"kepler-auth-go/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	cfg *config.Config
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{cfg: cfg}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check if user with same email exists in the same organization
	var existingUser models.User
	query := database.GetDB().Where("email = ?", req.Email)
	if req.OrganizationID != nil {
		query = query.Where("organization_id = ?", *req.OrganizationID)
	} else {
		query = query.Where("organization_id IS NULL")
	}

	if err := query.First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email already exists in this organization")
	}

	// Validate organization exists if provided
	if req.OrganizationID != nil {
		var org models.Organization
		if err := database.GetDB().First(&org, *req.OrganizationID).Error; err != nil {
			return nil, errors.New("organization not found")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:          req.Email,
		Name:           req.Name,
		Password:       string(hashedPassword),
		OrganizationID: req.OrganizationID,
		IsActive:       true,
	}

	if err := database.GetDB().Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	var user models.User
	query := database.GetDB().Preload("Groups").Preload("Organization").Where("email = ?", req.Email)

	// Filter by organization if provided
	if req.OrganizationID != nil {
		query = query.Where("organization_id = ?", *req.OrganizationID)
	} else {
		query = query.Where("organization_id IS NULL")
	}

	if err := query.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if user.IsDeleted || !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.generateToken(&user)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &models.LoginResponse{
		Token: token,
		User:  &user,
	}, nil
}

func (s *AuthService) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	var user models.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return database.GetDB().Model(&user).Update("password", string(hashedPassword)).Error
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
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

	claims := &middleware.Claims{
		UserID:         user.ID,
		Email:          user.Email,
		OrganizationID: user.OrganizationID,
		Permissions:    uniquePermissions,
		IsAdmin:        user.IsAdmin,
		IsStaff:        user.IsStaff,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.cfg.JWT.Expiration) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}
