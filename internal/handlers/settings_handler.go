package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"psclub-crm/internal/models"
	"psclub-crm/internal/services"
)

type SettingsHandler struct {
	service *services.SettingsService
}

func NewSettingsHandler(s *services.SettingsService) *SettingsHandler {
	return &SettingsHandler{service: s}
}

// GET /api/settings
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	set, err := h.service.GetSettings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, set)
}

// PUT /api/settings/:id
func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	var set models.Settings
	if err := c.ShouldBindJSON(&set); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.UpdateSettings(c.Request.Context(), &set)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, set)
}
