package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"psclub-crm/internal/models"
	"psclub-crm/internal/services"
)

type CashboxHandler struct {
	service *services.CashboxService
}

func NewCashboxHandlerCashboxHandler(service *services.CashboxService) *CashboxHandler {
	return &CashboxHandler{service: service}
}

// GET /api/cashbox
func (h *CashboxHandler) GetCashbox(c *gin.Context) {
	box, err := h.service.GetCashbox(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, box)
}

// PUT /api/cashbox/:id
func (h *CashboxHandler) UpdateCashbox(c *gin.Context) {
	var box models.Cashbox
	if err := c.ShouldBindJSON(&box); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.UpdateCashbox(c.Request.Context(), &box)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, box)
}
