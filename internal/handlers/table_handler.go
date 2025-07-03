package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"psclub-crm/internal/models"
	"psclub-crm/internal/services"
)

type TableHandler struct {
	service *services.TableService
}

func NewTableHandler(s *services.TableService) *TableHandler {
	return &TableHandler{service: s}
}

// POST /api/tables
func (h *TableHandler) CreateTable(c *gin.Context) {
	var table models.Table
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.service.CreateTable(c.Request.Context(), &table)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	table.ID = id
	c.JSON(http.StatusCreated, table)
}

// GET /api/tables
func (h *TableHandler) GetAllTables(c *gin.Context) {
	tables, err := h.service.GetAllTables(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tables)
}

// GET /api/tables/:id
func (h *TableHandler) GetTableByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	table, err := h.service.GetTableByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, table)
}

// PUT /api/tables/:id
func (h *TableHandler) UpdateTable(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var table models.Table
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	table.ID = id
	err = h.service.UpdateTable(id, &table)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, table)
}

// DELETE /api/tables/:id
func (h *TableHandler) DeleteTable(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err = h.service.DeleteTable(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// POST /api/tables/reorder
func (h *TableHandler) ReorderTable(c *gin.Context) {
	var req struct {
		ID     int `json:"id"`
		Number int `json:"number"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.ReorderTable(c.Request.Context(), req.ID, req.Number); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
