package handlers

import (
	"net/http"
	"strconv"
	"psclub-crm/internal/models"
	"psclub-crm/internal/services"
	"github.com/gin-gonic/gin"
)

type RepairHandler struct {
	service *services.RepairService
}

func NewRepairHandler(s *services.RepairService) *RepairHandler {
	return &RepairHandler{service: s}
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
	rep.ID = id
	err = h.service.UpdateRepair(c.Request.Context(), &rep)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rep)
}

// DELETE /api/repairs/:id
func (h *RepairHandler) DeleteRepair(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err = h.service.DeleteRepair(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
