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

type PricelistHistoryHandler struct {
	service  *services.PriceItemService
	expenses *services.ExpenseService
}

func NewPricelistHistoryHandler(s *services.PriceItemService, e *services.ExpenseService) *PricelistHistoryHandler {
	return &PricelistHistoryHandler{service: s, expenses: e}
}

// POST /api/pricelist-history
func (h *PricelistHistoryHandler) Create(c *gin.Context) {
	var hist models.PricelistHistory
	if err := c.ShouldBindJSON(&hist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	hist.CompanyID = companyID
	hist.BranchID = branchID
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	id, err := h.service.CreatePricelistHistory(ctx, &hist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	hist.ID = id
	c.JSON(http.StatusCreated, hist)
}

// GET /api/pricelist-history
func (h *PricelistHistoryHandler) GetAll(c *gin.Context) {
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	history, err := h.service.GetAllPricelistHistory(ctx)
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
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	history, err := h.service.GetPricelistHistoryByItem(ctx, id)
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
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	history, err := h.service.GetPricelistHistoryByCategory(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

func (h *PricelistHistoryHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	hist, err := h.service.GetPricelistHistoryByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.DeletePricelistHistory(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	item, _ := h.service.GetPriceItemByID(ctx, hist.PriceItemID)
	if item != nil {
		title := "Пополнение " + item.Name
		desc := "Пополнение товара " + item.Name + " в количестве " + strconv.FormatFloat(hist.Quantity, 'f', -1, 64) + " шт."
		_ = h.expenses.DeleteByDetails(c.Request.Context(), title, desc, hist.Total, 0)
	}
	c.Status(http.StatusNoContent)
}
