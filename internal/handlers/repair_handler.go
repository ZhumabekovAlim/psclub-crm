package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
	id, err := h.service.CreateRepair(c.Request.Context(), &rep)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rep.ID = id

	// ensure expense category "Ремонт" exists
	var catID int
	if cat, _ := h.expCats.GetByName(c.Request.Context(), "Ремонт"); cat != nil {
		catID = cat.ID
	} else {
		newCat := models.ExpenseCategory{Name: "Ремонт"}
		catID, _ = h.expCats.Create(c.Request.Context(), &newCat)
	}

	exp := models.Expense{
		Date:        time.Now(),
		Title:       "Починка, номер VIN: " + rep.VIN,
		Total:       rep.Price,
		Description: rep.Description,
		Paid:        false,
		CategoryID:  catID,
	}
	_, _ = h.expenses.CreateExpense(c.Request.Context(), &exp)

	c.JSON(http.StatusCreated, rep)
}

// GET /api/repairs
func (h *RepairHandler) GetAllRepairs(c *gin.Context) {
	repairs, err := h.service.GetAllRepairs(c.Request.Context())
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
	rep, err := h.service.GetRepairByID(c.Request.Context(), id)
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
	oldRep, _ := h.service.GetRepairByID(c.Request.Context(), id)
	rep.ID = id
	err = h.service.UpdateRepair(c.Request.Context(), &rep)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if oldRep != nil {
		oldTitle := "Починка, номер VIN: " + oldRep.VIN
		_ = h.expenses.DeleteByDetails(c.Request.Context(), oldTitle, oldRep.Description, oldRep.Price)
	}

	var catID int
	if cat, _ := h.expCats.GetByName(c.Request.Context(), "Ремонт"); cat != nil {
		catID = cat.ID
	} else {
		newCat := models.ExpenseCategory{Name: "Ремонт"}
		catID, _ = h.expCats.Create(c.Request.Context(), &newCat)
	}

	exp := models.Expense{
		Date:        time.Now(),
		Title:       "Починка, номер VIN: " + rep.VIN,
		Total:       rep.Price,
		Description: rep.Description,
		Paid:        false,
		CategoryID:  catID,
	}
	_, _ = h.expenses.CreateExpense(c.Request.Context(), &exp)

	c.JSON(http.StatusOK, rep)
}

// DELETE /api/repairs/:id
func (h *RepairHandler) DeleteRepair(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	rep, _ := h.service.GetRepairByID(c.Request.Context(), id)
	err = h.service.DeleteRepair(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rep != nil {
		title := "Починка, номер VIN: " + rep.VIN
		_ = h.expenses.DeleteByDetails(c.Request.Context(), title, rep.Description, rep.Price)
	}

	c.Status(http.StatusNoContent)
}
