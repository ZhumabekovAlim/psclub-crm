package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"psclub-crm/internal/models"
	"psclub-crm/internal/services"
	"strconv"
)

type EquipmentHandler struct {
	service       *services.EquipmentService
	inventoryServ *services.EquipmentInventoryService
}

func NewEquipmentHandler(s *services.EquipmentService, inv *services.EquipmentInventoryService) *EquipmentHandler {
	return &EquipmentHandler{service: s, inventoryServ: inv}
}

func (h *EquipmentHandler) Create(c *gin.Context) {
	var eq models.Equipment
	if err := c.ShouldBindJSON(&eq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	eq.CompanyID = c.GetInt("company_id")
	eq.BranchID = c.GetInt("branch_id")
	id, err := h.service.Create(c.Request.Context(), &eq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	eq.ID = id
	c.JSON(http.StatusCreated, eq)
}

func (h *EquipmentHandler) GetAll(c *gin.Context) {
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	list, err := h.service.GetAll(c.Request.Context(), companyID, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *EquipmentHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	eq, err := h.service.GetByID(c.Request.Context(), id, companyID, branchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, eq)
}

func (h *EquipmentHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var eq models.Equipment
	if err := c.ShouldBindJSON(&eq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	eq.ID = id
	eq.CompanyID = c.GetInt("company_id")
	eq.BranchID = c.GetInt("branch_id")
	if err := h.service.Update(c.Request.Context(), &eq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, eq)
}

func (h *EquipmentHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	if err := h.service.Delete(c.Request.Context(), id, companyID, branchID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *EquipmentHandler) PerformInventory(c *gin.Context) {
	var req struct {
		Items []services.EquipmentInventoryItem `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	if err := h.inventoryServ.PerformInventory(c.Request.Context(), req.Items, companyID, branchID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *EquipmentHandler) GetHistory(c *gin.Context) {
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	list, err := h.inventoryServ.GetHistory(c.Request.Context(), companyID, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
