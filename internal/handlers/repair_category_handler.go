package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"psclub-crm/internal/common"
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
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	cat.CompanyID = companyID
	cat.BranchID = branchID
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	id, err := h.service.Create(ctx, &cat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cat.ID = id
	c.JSON(http.StatusCreated, cat)
}

func (h *RepairCategoryHandler) GetAll(c *gin.Context) {
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	cats, err := h.service.GetAll(ctx)
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
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	cat.ID = id
	cat.CompanyID = companyID
	cat.BranchID = branchID
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	if err := h.service.Update(ctx, &cat); err != nil {
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
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	if err := h.service.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
