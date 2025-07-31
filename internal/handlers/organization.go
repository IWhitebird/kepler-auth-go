package handlers

import (
	"kepler-auth-go/internal/models"
	"kepler-auth-go/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	organizationService *services.OrganizationService
}

func NewOrganizationHandler() *OrganizationHandler {
	return &OrganizationHandler{
		organizationService: services.NewOrganizationService(),
	}
}

// GetOrganizations godoc
// @Summary Get all organizations with pagination and filtering
// @Description Get a paginated list of organizations with optional filtering (admin only)
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} models.PaginatedOrganizationResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/organizations [get]
func (h *OrganizationHandler) GetOrganizations(c *gin.Context) {
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

	response, err := h.organizationService.GetOrganizations(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetOrganization godoc
// @Summary Get organization by ID
// @Description Get a specific organization by their ID (admin only)
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Organization ID"
// @Success 200 {object} models.OrganizationResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/organizations/{id} [get]
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	response, err := h.organizationService.GetOrganizationByID(uint(id))
	if err != nil {
		if err.Error() == "organization not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateOrganization godoc
// @Summary Create a new organization
// @Description Create a new organization (admin only)
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.OrganizationCreateRequest true "Organization creation details"
// @Success 201 {object} models.OrganizationResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/organizations [post]
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req models.OrganizationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.organizationService.CreateOrganization(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// UpdateOrganization godoc
// @Summary Update organization by ID
// @Description Update a specific organization by their ID (admin only)
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Organization ID"
// @Param request body models.OrganizationUpdateRequest true "Organization update details"
// @Success 200 {object} models.OrganizationResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/organizations/{id} [patch]
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var req models.OrganizationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.organizationService.UpdateOrganization(uint(id), &req)
	if err != nil {
		if err.Error() == "organization not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteOrganization godoc
// @Summary Delete organization by ID
// @Description Delete a specific organization by their ID (admin only)
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Organization ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/organizations/{id} [delete]
func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	if err := h.organizationService.DeleteOrganization(uint(id)); err != nil {
		if err.Error() == "organization not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})
}
