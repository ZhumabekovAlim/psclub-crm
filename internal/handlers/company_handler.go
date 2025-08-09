package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"psclub-crm/internal/services"
)

type CompanyHandler struct {
	service *services.CompanyService
}

func NewCompanyHandler(s *services.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: s}
}

type createCompanyRequest struct {
	Name       string `json:"name"`
	BranchName string `json:"branch_name"`
}

func (h *CompanyHandler) Create(c *gin.Context) {
	var req createCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID, branchID, err := h.service.CreateCompany(c.Request.Context(), req.Name, req.BranchName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"company_id": companyID, "branch_id": branchID})
}
