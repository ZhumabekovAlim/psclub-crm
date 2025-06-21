package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"psclub-crm/internal/models"
	"psclub-crm/internal/services"
	"strconv"
	"time"
)

type PriceItemHandler struct {
	service  *services.PriceItemService
	expenses *services.ExpenseService
}

func NewPriceItemHandler(service *services.PriceItemService, expenseService *services.ExpenseService) *PriceItemHandler {
	return &PriceItemHandler{service: service, expenses: expenseService}
}

// CRUD
func (h *PriceItemHandler) CreatePriceItem(c *gin.Context) {
	var item models.PriceItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.service.CreatePriceItem(c.Request.Context(), &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	item.ID = id
	c.JSON(http.StatusCreated, item)
}

func (h *PriceItemHandler) GetAllPriceItems(c *gin.Context) {
	items, err := h.service.GetAllPriceItems(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *PriceItemHandler) GetPriceItemsByCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	items, err := h.service.GetPriceItemsByCategory(c.Request.Context(), id)
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
	item, err := h.service.GetPriceItemByID(c.Request.Context(), id)
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

	// Use pointers so that omitted JSON fields are nil and do not overwrite
	// existing values in the database.
	type updateInput struct {
		Name          *string  `json:"name"`
		CategoryID    *int     `json:"category_id"`
		SubcategoryID *int     `json:"subcategory_id"`
		Quantity      *int     `json:"quantity"`
		SalePrice     *float64 `json:"sale_price"`
		BuyPrice      *float64 `json:"buy_price"`
		IsSet         *bool    `json:"is_set"`
	}

	var in updateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.GetPriceItemByID(c.Request.Context(), id)
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
	err = h.service.UpdatePriceItem(c.Request.Context(), item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	err = h.service.DeletePriceItem(c.Request.Context(), id)
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
	hist.Operation = "INCOME"
	err := h.service.AddIncome(c.Request.Context(), &hist)
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
	hist.Operation = "OUTCOME"
	err := h.service.AddOutcome(c.Request.Context(), &hist)
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
	history, err := h.service.GetHistoryByItem(c.Request.Context(), itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

// Получить всю историю по всем товарам
func (h *PriceItemHandler) GetAllHistory(c *gin.Context) {
	history, err := h.service.GetAllHistory(c.Request.Context())
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
		Quantity int     `json:"quantity"`
		BuyPrice float64 `json:"buy_price"`
		UserID   int     `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hist := models.PricelistHistory{
		PriceItemID: itemID,
		Quantity:    in.Quantity,
		BuyPrice:    in.BuyPrice,
		Total:       in.BuyPrice * float64(in.Quantity),
		UserID:      in.UserID,
	}
	if err := h.service.Replenish(c.Request.Context(), &hist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create corresponding expense entry
	item, _ := h.service.GetPriceItemByID(c.Request.Context(), itemID)
	exp := models.Expense{
		Date:       time.Now(),
		Title:      "Replenish " + item.Name,
		Total:      hist.Total,
		Paid:       false,
		CategoryID: 0,
	}
	_, _ = h.expenses.CreateExpense(c.Request.Context(), &exp)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
