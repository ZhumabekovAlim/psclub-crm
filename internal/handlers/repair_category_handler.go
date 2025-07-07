package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"psclub-crm/internal/models"
	"psclub-crm/internal/services"
)

type RepairCategoryHandler struct {
	service *services.RepairCategoryService
}

func NewRepairCategoryHandler(s *services.RepairCategoryService) *RepairCategoryHandler {
	return &RepairCategoryHandler{service: s}
}

func (h *RepairCategoryHandler) Create(c *gin.Context) {
	var cat models.RepairCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.service.Create(c.Request.Context(), &cat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cat.ID = id
	c.JSON(http.StatusCreated, cat)
}

func (h *RepairCategoryHandler) GetAll(c *gin.Context) {
	cats, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cats)
}

func (h *RepairCategoryHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var cat models.RepairCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat.ID = id
	if err := h.service.Update(c.Request.Context(), &cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *RepairCategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
