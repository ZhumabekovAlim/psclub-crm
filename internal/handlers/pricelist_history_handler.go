package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"psclub-crm/internal/models"
	"psclub-crm/internal/services"
	"strconv"
)

type PricelistHistoryHandler struct {
	service *services.PriceItemService
}

func NewPricelistHistoryHandler(s *services.PriceItemService) *PricelistHistoryHandler {
	return &PricelistHistoryHandler{service: s}
}

// POST /api/pricelist-history
func (h *PricelistHistoryHandler) Create(c *gin.Context) {
	var hist models.PricelistHistory
	if err := c.ShouldBindJSON(&hist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.service.CreatePricelistHistory(c.Request.Context(), &hist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	hist.ID = id
	c.JSON(http.StatusCreated, hist)
}

// GET /api/pricelist-history
func (h *PricelistHistoryHandler) GetAll(c *gin.Context) {
	history, err := h.service.GetAllPricelistHistory(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

// GET /api/pricelist-history/item/:id
func (h *PricelistHistoryHandler) GetByItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	history, err := h.service.GetPricelistHistoryByItem(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

func (h *PricelistHistoryHandler) GetByCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	history, err := h.service.GetPricelistHistoryByCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}
