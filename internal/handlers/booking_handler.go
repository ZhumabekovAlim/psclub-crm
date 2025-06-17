package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
	id, err := h.service.CreateBooking(c.Request.Context(), &b)
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
	bookings, err := h.service.GetAllBookings(c.Request.Context())
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
	b, err := h.service.GetBookingByID(c.Request.Context(), id)
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
	err = h.service.UpdateBooking(c.Request.Context(), &b)
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
	err = h.service.DeleteBooking(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
