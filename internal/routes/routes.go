package routes

import (
	"go-rest/internal/handlers"
	"go-rest/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Public routes
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		// Protected routes
		items := api.Group("/items")
		items.Use(middleware.AuthMiddleware())
		{
			items.POST("", handlers.CreateItem)
			items.GET("", handlers.GetItems)
			items.GET("/:id", handlers.GetItem)
			items.PUT("/:id", handlers.UpdateItem)
			items.DELETE("/:id", handlers.DeleteItem)
		}
	}
}
