package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"psclub-crm/internal/models"
	"psclub-crm/internal/services"
)

type PaymentTypeHandler struct {
	service *services.PaymentTypeService
}

func NewPaymentTypeHandler(s *services.PaymentTypeService) *PaymentTypeHandler {
	return &PaymentTypeHandler{service: s}
}

// POST /api/payment-types
func (h *PaymentTypeHandler) CreatePaymentType(c *gin.Context) {
	var pt models.PaymentType
	if err := c.ShouldBindJSON(&pt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.service.CreatePaymentType(c.Request.Context(), &pt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	pt.ID = id
	c.JSON(http.StatusCreated, pt)
}

// GET /api/payment-types
func (h *PaymentTypeHandler) GetAllPaymentTypes(c *gin.Context) {
	pts, err := h.service.GetAllPaymentTypes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pts)
}

// PUT /api/payment-types/:id
func (h *PaymentTypeHandler) UpdatePaymentType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var pt models.PaymentType
	if err := c.ShouldBindJSON(&pt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pt.ID = id
	if err := h.service.UpdatePaymentType(c.Request.Context(), &pt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pt)
}

// DELETE /api/payment-types/:id
func (h *PaymentTypeHandler) DeletePaymentType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.DeletePaymentType(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
