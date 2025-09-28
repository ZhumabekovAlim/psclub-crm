package handlers

import (
	"database/sql"
	"errors"
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
	subcategory.CompanyID = c.GetInt("company_id")
	subcategory.BranchID = c.GetInt("branch_id")
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
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	subcategories, err := h.service.GetAllSubcategories(c.Request.Context(), companyID, branchID)
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
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	subcategory, err := h.service.GetSubcategoryByID(c.Request.Context(), id, companyID, branchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subcategory)
}

func (h *SubcategoryHandler) GetSubcategoryByCategoryID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	subcategory, err := h.service.GetSubcategoryByCategoryID(c.Request.Context(), id, companyID, branchID)
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
	var input struct {
		CategoryID *int    `json:"category_id"`
		Name       *string `json:"name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.CategoryID == nil && input.Name == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	subcategory, err := h.service.GetSubcategoryByID(c.Request.Context(), id, companyID, branchID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "subcategory not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if input.CategoryID != nil {
		subcategory.CategoryID = *input.CategoryID
	}
	if input.Name != nil {
		subcategory.Name = *input.Name
	}
	err = h.service.UpdateSubcategory(c.Request.Context(), subcategory)
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
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	err = h.service.DeleteSubcategory(c.Request.Context(), id, companyID, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
