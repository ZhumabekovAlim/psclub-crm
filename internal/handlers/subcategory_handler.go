package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"psclub-crm/internal/models"
	"psclub-crm/internal/services"
)

type SubcategoryHandler struct {
	service *services.SubcategoryService
}

func NewSubcategoryHandler(s *services.SubcategoryService) *SubcategoryHandler {
	return &SubcategoryHandler{service: s}
}

// POST /api/subcategories
func (h *SubcategoryHandler) CreateSubcategory(c *gin.Context) {
	var subcategory models.Subcategory
	if err := c.ShouldBindJSON(&subcategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.service.CreateSubcategory(c.Request.Context(), &subcategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	subcategory.ID = id
	c.JSON(http.StatusCreated, subcategory)
}

// GET /api/subcategories
func (h *SubcategoryHandler) GetAllSubcategories(c *gin.Context) {
	subcategories, err := h.service.GetAllSubcategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subcategories)
}

// GET /api/subcategories/:id
func (h *SubcategoryHandler) GetSubcategoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	subcategory, err := h.service.GetSubcategoryByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subcategory)
}

// PUT /api/subcategories/:id
func (h *SubcategoryHandler) UpdateSubcategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var subcategory models.Subcategory
	if err := c.ShouldBindJSON(&subcategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	subcategory.ID = id
	err = h.service.UpdateSubcategory(c.Request.Context(), &subcategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subcategory)
}

// DELETE /api/subcategories/:id
func (h *SubcategoryHandler) DeleteSubcategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err = h.service.DeleteSubcategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
