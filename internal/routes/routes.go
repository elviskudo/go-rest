package routes

import (
	"go-rest/internal/handlers"
	"go-rest/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

			// Product Enhancements
			items.POST("/:id/media", handlers.UploadItemMedia)
			items.POST("/:id/reviews", handlers.CreateReview)
			items.POST("/:id/favorite", handlers.ToggleFavorite)
		}

		// Categories
		categories := api.Group("/categories")
		categories.Use(middleware.AuthMiddleware())
		{
			categories.POST("", handlers.CreateCategory)
			categories.GET("", handlers.GetCategories)
		}

		// Warehouses
		warehouses := api.Group("/warehouses")
		warehouses.Use(middleware.AuthMiddleware())
		{
			warehouses.POST("", handlers.CreateWarehouse)
			warehouses.GET("", handlers.GetWarehouses)
		}

		// Suppliers
		suppliers := api.Group("/suppliers")
		suppliers.Use(middleware.AuthMiddleware())
		{
			suppliers.POST("", handlers.CreateSupplier)
			suppliers.GET("", handlers.GetSuppliers)
		}

		// Discounts
		discounts := api.Group("/discounts")
		discounts.Use(middleware.AuthMiddleware())
		{
			discounts.POST("", handlers.CreateDiscount)
			discounts.GET("", handlers.GetDiscounts)
		}

		// Inventory
		inventory := api.Group("/inventory")
		inventory.Use(middleware.AuthMiddleware())
		{
			inventory.POST("/add", handlers.AddStock)
			inventory.POST("/transfer", handlers.TransferStock)
		}

		// Purchase Orders
		pos := api.Group("/purchase-orders")
		pos.Use(middleware.AuthMiddleware())
		{
			pos.POST("", handlers.CreatePurchaseOrder)
			pos.PUT("/:id/status", handlers.UpdatePurchaseOrderStatus)
		}

		// Sales Orders
		orders := api.Group("/orders")
		orders.Use(middleware.AuthMiddleware())
		{
			orders.POST("", handlers.CreateOrder)
		}

		// Reports & Dashboard
		reports := api.Group("/reports")
		reports.Use(middleware.AuthMiddleware())
		{
			reports.GET("/financial", handlers.GetFinancialReport)
			reports.GET("/sales", handlers.GetSalesReport)
			reports.GET("/dashboard", handlers.GetDashboardSummary)
		}
	}

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
