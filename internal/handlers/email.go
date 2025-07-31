package handlers

import (
	"kepler-auth-go/internal/config"
	"kepler-auth-go/internal/models"
	"kepler-auth-go/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	emailService *services.EmailService
}

func NewEmailHandler(cfg *config.Config) *EmailHandler {
	return &EmailHandler{
		emailService: services.NewEmailService(cfg),
	}
}

// SendEmail godoc
// @Summary Send email
// @Description Send email to specified recipient
// @Tags email
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.EmailRequest true "Email details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/email/send [post]
func (h *EmailHandler) SendEmail(c *gin.Context) {
	var req models.EmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.emailService.SendEmail(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}
