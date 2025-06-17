package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"psclub-crm/internal/models"
	"psclub-crm/internal/services"
)

type PriceSetHandler struct {
	service *services.PriceSetService
}

func NewPriceSetHandler(s *services.PriceSetService) *PriceSetHandler {
	return &PriceSetHandler{service: s}
}

func (h *PriceSetHandler) CreatePriceSet(c *gin.Context) {
	var ps models.PriceSet
	if err := c.ShouldBindJSON(&ps); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.service.CreatePriceSet(c.Request.Context(), &ps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ps.ID = id
	c.JSON(http.StatusCreated, ps)
}

func (h *PriceSetHandler) GetAllPriceSets(c *gin.Context) {
	sets, err := h.service.GetAllPriceSets(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sets)
}

func (h *PriceSetHandler) GetPriceSetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	set, err := h.service.GetPriceSetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, set)
}

func (h *PriceSetHandler) UpdatePriceSet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var ps models.PriceSet
	if err := c.ShouldBindJSON(&ps); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ps.ID = id
	if err := h.service.UpdatePriceSet(c.Request.Context(), &ps); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ps)
}

func (h *PriceSetHandler) DeletePriceSet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.DeletePriceSet(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
