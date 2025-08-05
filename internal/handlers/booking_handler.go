package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"psclub-crm/internal/common"
	"psclub-crm/internal/models"
	"psclub-crm/internal/services"
	"strconv"
)

type BookingHandler struct {
	service *services.BookingService
}

func NewBookingHandler(s *services.BookingService) *BookingHandler {
	return &BookingHandler{service: s}
}

// POST /api/bookings
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var b models.Booking
	if err := c.ShouldBindJSON(&b); err != nil {
		log.Printf("create booking bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, it := range b.Items {
		if it.Quantity < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "quantity cannot be negative"})

		}
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	id, err := h.service.CreateBooking(ctx, &b)
	if err != nil {
		log.Printf("create booking service error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	b.ID = id
	c.JSON(http.StatusCreated, b)
}

// GET /api/bookings
func (h *BookingHandler) GetAllBookings(c *gin.Context) {
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	bookings, err := h.service.GetAllBookings(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// GET /api/bookings/client/:id
func (h *BookingHandler) GetBookingsByClientID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	bookings, err := h.service.GetBookingsByClientID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// GET /api/bookings/:id
func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	b, err := h.service.GetBookingByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

// PUT /api/bookings/:id
func (h *BookingHandler) UpdateBooking(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var b models.Booking
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	b.ID = id
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	err = h.service.UpdateBooking(ctx, &b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

// DELETE /api/bookings/:id
func (h *BookingHandler) DeleteBooking(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	err = h.service.DeleteBooking(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
