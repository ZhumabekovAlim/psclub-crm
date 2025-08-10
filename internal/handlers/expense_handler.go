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

type ExpenseHandler struct {
	service *services.ExpenseService
}

func NewExpenseHandler(s *services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: s}
}

// POST /api/expenses
func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var e models.Expense
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	e.CompanyID = companyID
	e.BranchID = branchID
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	id, err := h.service.CreateExpense(ctx, &e)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	e.ID = id
	c.JSON(http.StatusCreated, e)
}

// GET /api/expenses
func (h *ExpenseHandler) GetAllExpenses(c *gin.Context) {
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)

	expenses, err := h.service.GetAllExpenses(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}

// GET /api/expenses/:id
func (h *ExpenseHandler) GetExpenseByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	e, err := h.service.GetExpenseByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, e)
}

// PUT /api/expenses/:id
func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var e models.Expense
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	e.ID = id
	e.CompanyID = companyID
	e.BranchID = branchID
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	err = h.service.UpdateExpense(ctx, &e)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, e)
}

// DELETE /api/expenses/:id
func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	err = h.service.DeleteExpense(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
