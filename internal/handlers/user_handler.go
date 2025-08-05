package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"psclub-crm/internal/models"
	"psclub-crm/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// POST /api/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.CompanyID = c.GetInt("company_id")
	user.BranchID = c.GetInt("branch_id")

	id, err := h.service.CreateUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.ID = id
	c.JSON(http.StatusCreated, user)
}

// GET /api/users
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	users, err := h.service.GetAllUsers(c.Request.Context(), companyID, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Пароли не выдаём!
	for i := range users {
		users[i].Password = ""
	}
	c.JSON(http.StatusOK, users)
}

// GET /api/users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	user, err := h.service.GetUserByID(c.Request.Context(), id, companyID, branchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// PUT /api/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = id
	user.CompanyID = c.GetInt("company_id")
	user.BranchID = c.GetInt("branch_id")
	err = h.service.UpdateUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// DELETE /api/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	companyID := c.GetInt("company_id")
	branchID := c.GetInt("branch_id")
	err = h.service.DeleteUser(c.Request.Context(), id, companyID, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
