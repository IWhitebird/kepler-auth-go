package handlers

import (
	"kepler-auth-go/internal/models"
	"kepler-auth-go/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
	}
}

// GetUsers godoc
// @Summary Get all users with pagination and filtering
// @Description Get a paginated list of users with optional filtering
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Param status query string false "User status filter" Enums(active, pending, deactivated)
// @Param is_active query bool false "Active status filter"
// @Success 200 {object} models.PaginatedUserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	query := &models.PaginationQuery{
		Page:     1,
		PageSize: 10,
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			query.Page = p
		}
	}

	if pageSize := c.Query("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
			query.PageSize = ps
		}
	}

	query.Search = c.Query("search")
	query.Status = c.Query("status")

	if isActive := c.Query("is_active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			query.IsActive = &active
		}
	}

	// Get organization context from JWT claims
	organizationID, _ := c.Get("organization_id")
	var orgID *uint
	if organizationID != nil {
		if oid, ok := organizationID.(*uint); ok {
			orgID = oid
		}
	}

	response, err := h.userService.GetUsers(query, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get a specific user by their ID
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get organization context from JWT claims
	organizationID, _ := c.Get("organization_id")
	var orgID *uint
	if organizationID != nil {
		if oid, ok := organizationID.(*uint); ok {
			orgID = oid
		}
	}

	response, err := h.userService.GetUserByID(uint(id), orgID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUser godoc
// @Summary Update user by ID
// @Description Update a specific user by their ID (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body models.UserUpdateRequest true "User update details"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/users/{id} [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get organization context from JWT claims
	organizationID, _ := c.Get("organization_id")
	var orgID *uint
	if organizationID != nil {
		if oid, ok := organizationID.(*uint); ok {
			orgID = oid
		}
	}

	response, err := h.userService.UpdateUser(uint(id), &req, orgID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteUser godoc
// @Summary Delete user by ID
// @Description Soft delete a specific user by their ID (admin only)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get organization context from JWT claims
	organizationID, _ := c.Get("organization_id")
	var orgID *uint
	if organizationID != nil {
		if oid, ok := organizationID.(*uint); ok {
			orgID = oid
		}
	}

	if err := h.userService.DeleteUser(uint(id), orgID); err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
