package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"psclub-crm/internal/common"
	"psclub-crm/internal/models"
	"psclub-crm/internal/services"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(c)
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	u.CompanyID = companyID
	u.BranchID = branchID
	ctx := context.WithValue(c.Request.Context(), common.CtxCompanyID, companyID)
	ctx = context.WithValue(ctx, common.CtxBranchID, branchID)
	access, refresh, err := h.service.Register(ctx, &u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": u, "access": access, "refresh": refresh})
}

// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	access, refresh, role, permission, name, id, companyID, branchID, err := h.service.Login(c.Request.Context(), req.Phone, req.Password)
	if err != nil {
		log.Printf("login error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access": access, "refresh": refresh, "id": id, "role": role, "permissions": permission, "name": name, "company_id": companyID, "branch_id": branchID})
}

// POST /api/auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		Refresh string `json:"refresh"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	access, refresh, err := h.service.Refresh(c.Request.Context(), req.Refresh)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access": access, "refresh": refresh})
}
