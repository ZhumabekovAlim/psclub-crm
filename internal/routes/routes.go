package internal

import (
	"github.com/gin-gonic/gin"
	"psclub-crm/internal/handlers"
	"psclub-crm/internal/services"
)

func InitRoutes(r *gin.Engine, clientService *services.ClientService) {
	clientHandler := handlers.NewClientHandler(clientService)

	api := r.Group("/api")
	{
		clients := api.Group("/clients")
		{
			clients.POST("", clientHandler.CreateClient)
			clients.GET("", clientHandler.GetAllClients)
			clients.GET("/:id", clientHandler.GetClientByID)
			clients.PUT("/:id", clientHandler.UpdateClient)
			clients.DELETE("/:id", clientHandler.DeleteClient)
		}
	}
}
