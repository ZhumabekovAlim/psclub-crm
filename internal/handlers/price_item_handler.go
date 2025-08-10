package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"psclub-crm/internal/common"
	"psclub-crm/internal/models"
	"psclub-crm/internal/services"
)

type PriceItemHandler struct {
	service    *services.PriceItemService
	expenses   *services.ExpenseService
	expCats    *services.ExpenseCategoryService
	categories *services.CategoryService
}

func NewPriceItemHandler(service *services.PriceItemService, expenseService *services.ExpenseService, expCatService *services.ExpenseCategoryService, categoryService *services.CategoryService) *PriceItemHandler {
	return &PriceItemHandler{service: service, expenses: expenseService, expCats: expCatService, categories: categoryService}
}

// CRUD
func (h *PriceItemHandler) CreatePriceItem(c *gin.Context) {
	var item models.PriceItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	item.CompanyID = companyID
	item.BranchID = branchID
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	id, err := h.service.CreatePriceItem(ctx, &item)
	if err != nil {
		if err == services.ErrNameExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	item.ID = id
	c.JSON(http.StatusCreated, item)
}

func (h *PriceItemHandler) GetAllPriceItems(c *gin.Context) {
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	items, err := h.service.GetAllPriceItems(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *PriceItemHandler) GetPriceItemsByCategoryName(c *gin.Context) {
	name := strings.TrimSpace(c.Param("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category name required"})
		return
	}

	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")

	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)

	items, err := h.service.GetPriceItemsByCategoryName(ctx, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *PriceItemHandler) GetPriceItemByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	item, err := h.service.GetPriceItemByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *PriceItemHandler) UpdatePriceItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	// Use pointers so that omitted JSON fields are nil and do not overwrite
	// existing values in the database.
	type updateInput struct {
		Name          *string  `json:"name"`
		CategoryID    *int     `json:"category_id"`
		SubcategoryID *int     `json:"subcategory_id"`
		Quantity      *float64 `json:"quantity"`
		SalePrice     *float64 `json:"sale_price"`
		BuyPrice      *float64 `json:"buy_price"`
		IsSet         *bool    `json:"is_set"`
	}

	var in updateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.GetPriceItemByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if in.Name != nil {
		item.Name = *in.Name
	}
	if in.CategoryID != nil {
		item.CategoryID = *in.CategoryID
	}
	if in.SubcategoryID != nil {
		item.SubcategoryID = *in.SubcategoryID
	}
	if in.Quantity != nil {
		item.Quantity = *in.Quantity
	}
	if in.SalePrice != nil {
		item.SalePrice = *in.SalePrice
	}
	if in.BuyPrice != nil {
		item.BuyPrice = *in.BuyPrice
	}
	if in.IsSet != nil {
		item.IsSet = *in.IsSet
	}

	item.ID = id
	err = h.service.UpdatePriceItem(ctx, item)
	if err != nil {
		if err == services.ErrNameExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *PriceItemHandler) DeletePriceItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	err = h.service.DeletePriceItem(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// Добавить приход (INCOME)
func (h *PriceItemHandler) AddIncome(c *gin.Context) {
	var hist models.PriceItemHistory
	if err := c.ShouldBindJSON(&hist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	hist.Operation = "INCOME"
	err := h.service.AddIncome(ctx, &hist)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// Добавить расход (OUTCOME)
func (h *PriceItemHandler) AddOutcome(c *gin.Context) {
	var hist models.PriceItemHistory
	if err := c.ShouldBindJSON(&hist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	hist.Operation = "OUTCOME"
	err := h.service.AddOutcome(ctx, &hist)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// Получить всю историю по товару
func (h *PriceItemHandler) GetHistoryByItem(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price_item_id"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	history, err := h.service.GetHistoryByItem(ctx, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

// Получить всю историю по всем товарам
func (h *PriceItemHandler) GetAllHistory(c *gin.Context) {
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	history, err := h.service.GetAllHistory(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

// Replenish stock for a price item
func (h *PriceItemHandler) Replenish(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var in struct {
		Quantity float64 `json:"quantity"`
		BuyPrice float64 `json:"buy_price"`
		UserID   int     `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	hist := models.PricelistHistory{
		PriceItemID: itemID,
		Quantity:    in.Quantity,
		BuyPrice:    in.BuyPrice,
		Total:       in.BuyPrice * in.Quantity,
		UserID:      in.UserID,
	}
	if err := h.service.Replenish(ctx, &hist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create corresponding expense entry with automatic category mapping
	item, _ := h.service.GetPriceItemByID(ctx, itemID)
	cat, _ := h.categories.GetCategoryByID(c.Request.Context(), item.CategoryID, companyID, branchID)
	var expCatID int
	if cat != nil {
		if ec, _ := h.expCats.GetByName(ctx, cat.Name); ec != nil {
			expCatID = ec.ID
		} else {
			newCat := models.ExpenseCategory{Name: cat.Name, CompanyID: companyID, BranchID: branchID}
			expCatID, _ = h.expCats.Create(ctx, &newCat)
		}
	}

	exp := models.Expense{
		Date:        time.Now(),
		Title:       "Пополнение " + item.Name,
		Total:       hist.Total,
		Paid:        false,
		CategoryID:  expCatID,
		Description: "Пополнение товара " + item.Name + " в количестве " + strconv.FormatFloat(in.Quantity, 'f', -1, 64) + " шт.",
	}
	_, _ = h.expenses.CreateExpense(ctx, &exp)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
