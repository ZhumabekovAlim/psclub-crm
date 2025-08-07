package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"psclub-crm/internal/common"
	"psclub-crm/internal/models"
	"psclub-crm/internal/services"
	"strconv"
	"time"
)

type RepairHandler struct {
	service  *services.RepairService
	expenses *services.ExpenseService
	expCats  *services.ExpenseCategoryService
}

func NewRepairHandler(s *services.RepairService, expenseService *services.ExpenseService, catService *services.ExpenseCategoryService) *RepairHandler {
	return &RepairHandler{service: s, expenses: expenseService, expCats: catService}
}

// POST /api/repairs
func (h *RepairHandler) CreateRepair(c *gin.Context) {
	var rep models.Repair
	if err := c.ShouldBindJSON(&rep); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	id, err := h.service.CreateRepair(ctx, &rep)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rep.ID = id

	// ensure expense category "Ремонт" exists
	var catID int
	if cat, _ := h.expCats.GetByName(ctx, "Ремонт"); cat != nil {
		catID = cat.ID
	} else {
		newCat := models.ExpenseCategory{Name: "Ремонт", CompanyID: companyID, BranchID: branchID}
		catID, _ = h.expCats.Create(ctx, &newCat)
	}

	exp := models.Expense{
		Date:             time.Now(),
		Title:            "Починка, номер VIN: " + rep.VIN,
		Total:            rep.Price,
		Description:      rep.Description,
		Paid:             false,
		CategoryID:       catID,
		RepairCategoryID: rep.CategoryID,
	}
	_, _ = h.expenses.CreateExpense(ctx, &exp)

	c.JSON(http.StatusCreated, rep)
}

// GET /api/repairs
func (h *RepairHandler) GetAllRepairs(c *gin.Context) {
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	repairs, err := h.service.GetAllRepairs(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, repairs)
}

// GET /api/repairs/:id
func (h *RepairHandler) GetRepairByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	rep, err := h.service.GetRepairByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rep)
}

// PUT /api/repairs/:id
func (h *RepairHandler) UpdateRepair(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var rep models.Repair
	if err := c.ShouldBindJSON(&rep); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	oldRep, _ := h.service.GetRepairByID(ctx, id)
	rep.ID = id
	err = h.service.UpdateRepair(ctx, &rep)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if oldRep != nil {
		oldTitle := "Починка, номер VIN: " + oldRep.VIN
		_ = h.expenses.DeleteByDetails(ctx, oldTitle, oldRep.Description, oldRep.Price, oldRep.CategoryID)
	}

	var catID int
	if cat, _ := h.expCats.GetByName(ctx, "Ремонт"); cat != nil {
		catID = cat.ID
	} else {
		newCat := models.ExpenseCategory{Name: "Ремонт", CompanyID: companyID, BranchID: branchID}
		catID, _ = h.expCats.Create(ctx, &newCat)
	}

	exp := models.Expense{
		Date:             time.Now(),
		Title:            "Починка, номер VIN: " + rep.VIN,
		Total:            rep.Price,
		Description:      rep.Description,
		Paid:             false,
		CategoryID:       catID,
		RepairCategoryID: rep.CategoryID,
	}
	_, _ = h.expenses.CreateExpense(ctx, &exp)

	c.JSON(http.StatusOK, rep)
}

// DELETE /api/repairs/:id
func (h *RepairHandler) DeleteRepair(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	rep, _ := h.service.GetRepairByID(ctx, id)
	err = h.service.DeleteRepair(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rep != nil {
		title := "Починка, номер VIN: " + rep.VIN
		_ = h.expenses.DeleteByDetails(ctx, title, rep.Description, rep.Price, rep.CategoryID)
	}

	c.Status(http.StatusNoContent)
}
