package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"psclub-crm/internal/services"
)

type InventoryHandler struct {
	service *services.InventoryService
}

func NewInventoryHandler(s *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: s}
}

func (h *InventoryHandler) PerformInventory(c *gin.Context) {
	var req struct {
		Items []services.InventoryItem `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.PerformInventory(c.Request.Context(), req.Items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *InventoryHandler) GetHistory(c *gin.Context) {
	history, err := h.service.GetHistory(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}
