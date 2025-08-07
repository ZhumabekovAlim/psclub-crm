package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"psclub-crm/internal/models"
	"psclub-crm/internal/services"
	"strconv"
)

type SettingsHandler struct {
	service *services.SettingsService
}

func NewSettingsHandler(s *services.SettingsService) *SettingsHandler {
	return &SettingsHandler{service: s}
}

// GET /api/settings
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	set, err := h.service.GetSettings(c.Request.Context(), companyID, branchID)
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
	set.CompanyID = c.GetInt("company_id")
	set.BranchID = c.GetInt("branch_id")
	err := h.service.UpdateSettings(c.Request.Context(), &set)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, set)
}

// POST /api/settings
func (h *SettingsHandler) CreateSettings(c *gin.Context) {
	var set models.Settings
	if err := c.ShouldBindJSON(&set); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	set.CompanyID = c.GetInt("company_id")
	set.BranchID = c.GetInt("branch_id")
	id, err := h.service.CreateSettings(c.Request.Context(), &set)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	set.ID = id
	c.JSON(http.StatusCreated, set)
}

// DELETE /api/settings/:id
func (h *SettingsHandler) DeleteSettings(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	if err := h.service.DeleteSettings(c.Request.Context(), id, companyID, branchID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *SettingsHandler) GetTablesCount(c *gin.Context) {
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	cnt, err := h.service.GetTablesCount(c.Request.Context(), companyID, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tables_count": cnt})
}

func (h *SettingsHandler) GetNotificationTime(c *gin.Context) {
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	n, err := h.service.GetNotificationTime(c.Request.Context(), companyID, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notification_time": n})
}
