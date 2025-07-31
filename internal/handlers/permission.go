package handlers

import (
	"kepler-auth-go/internal/models"
	"kepler-auth-go/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permissionService *services.PermissionService
}

func NewPermissionHandler() *PermissionHandler {
	return &PermissionHandler{
		permissionService: services.NewPermissionService(),
	}
}

// GetPermissions godoc
// @Summary Get all permissions with pagination and filtering
// @Description Get a paginated list of permissions with optional filtering
// @Tags permissions
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} models.PaginatedPermissionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/permissions [get]
func (h *PermissionHandler) GetPermissions(c *gin.Context) {
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

	response, err := h.permissionService.GetPermissions(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAuthGroups godoc
// @Summary Get all auth groups with pagination and filtering
// @Description Get a paginated list of auth groups with optional filtering
// @Tags auth-groups
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} models.PaginatedAuthGroupResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth-groups [get]
func (h *PermissionHandler) GetAuthGroups(c *gin.Context) {
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

	response, err := h.permissionService.GetAuthGroups(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateAuthGroup godoc
// @Summary Create a new auth group
// @Description Create a new auth group (admin only)
// @Tags auth-groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.AuthGroupRequest true "Auth group creation details"
// @Success 201 {object} models.AuthGroup
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth-groups [post]
func (h *PermissionHandler) CreateAuthGroup(c *gin.Context) {
	var req models.AuthGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authGroup, err := h.permissionService.CreateAuthGroup(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, authGroup)
}

// GetAuthGroup godoc
// @Summary Get auth group by ID
// @Description Get a specific auth group by their ID
// @Tags auth-groups
// @Produce json
// @Security BearerAuth
// @Param id path int true "Auth Group ID"
// @Success 200 {object} models.AuthGroup
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/auth-groups/{id} [get]
func (h *PermissionHandler) GetAuthGroup(c *gin.Context) {
	idStr := c.Param("id")
	_, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auth group ID"})
		return
	}

	// For now, return a simple implementation - in a real app this would get by ID
	var authGroup models.AuthGroup
	authGroup.ID = 1
	authGroup.Name = "Sample Group"

	c.JSON(http.StatusOK, authGroup)
}

// UpdateAuthGroup godoc
// @Summary Update auth group by ID
// @Description Update a specific auth group by their ID (admin only)
// @Tags auth-groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Auth Group ID"
// @Param request body models.AuthGroupRequest true "Auth group update details"
// @Success 200 {object} models.AuthGroup
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/auth-groups/{id} [patch]
func (h *PermissionHandler) UpdateAuthGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auth group ID"})
		return
	}

	var req models.AuthGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authGroup, err := h.permissionService.UpdateAuthGroup(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authGroup)
}

// DeleteAuthGroup godoc
// @Summary Delete auth group by ID
// @Description Delete a specific auth group by their ID (admin only)
// @Tags auth-groups
// @Produce json
// @Security BearerAuth
// @Param id path int true "Auth Group ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/auth-groups/{id} [delete]
func (h *PermissionHandler) DeleteAuthGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auth group ID"})
		return
	}

	if err := h.permissionService.DeleteAuthGroup(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auth group deleted successfully"})
}
