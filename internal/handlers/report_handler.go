package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"psclub-crm/internal/services"
	"strconv"
	"time"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) GetSummaryReport(c *gin.Context) {
	from, to, tFrom, tTo := getPeriod(c)
	userID := getUserID(c)
	data, err := h.service.SummaryReport(c.Request.Context(), from, to, tFrom, tTo, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *ReportHandler) GetAdminsReport(c *gin.Context) {
	from, to, tFrom, tTo := getPeriod(c)
	userID := getUserID(c)
	data, err := h.service.AdminsReport(c.Request.Context(), from, to, tFrom, tTo, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *ReportHandler) GetSalesReport(c *gin.Context) {
	from, to, tFrom, tTo := getPeriod(c)
	userID := getUserID(c)
	data, err := h.service.SalesReport(c.Request.Context(), from, to, tFrom, tTo, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *ReportHandler) GetAnalyticsReport(c *gin.Context) {
	from, to, tFrom, tTo := getPeriod(c)
	userID := getUserID(c)
	data, err := h.service.AnalyticsReport(c.Request.Context(), from, to, tFrom, tTo, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *ReportHandler) GetDiscountsReport(c *gin.Context) {
	from, to, tFrom, tTo := getPeriod(c)
	userID := getUserID(c)
	data, err := h.service.DiscountsReport(c.Request.Context(), from, to, tFrom, tTo, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func getPeriod(c *gin.Context) (from, to time.Time, tFrom, tTo string) {
	layoutDate := "2006-01-02"
	fromStr := c.DefaultQuery("from", time.Now().AddDate(0, 0, -7).Format(layoutDate))
	toStr := c.DefaultQuery("to", time.Now().Format(layoutDate))
	tFrom = c.DefaultQuery("tFrom", "00:00")
	tTo = c.DefaultQuery("tTo", "23:59:59")
	from, _ = time.Parse(layoutDate, fromStr)
	to, _ = time.Parse(layoutDate, toStr)
	return from, to, tFrom, tTo
}

func getUserID(c *gin.Context) int {
	userStr := c.DefaultQuery("user_id", "all")
	if userStr == "all" {
		return 0
	}
	id, err := strconv.Atoi(userStr)
	if err != nil {
		return 0
	}
	return id
}
