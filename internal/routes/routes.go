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
		// Items
		items := api.Group("/items")
		items.Use(middleware.AuthMiddleware())
		{
			items.POST("", middleware.RequirePermission("items", "write"), handlers.CreateItem)
			items.GET("", middleware.RequirePermission("items", "read"), handlers.GetItems)
			items.GET("/:id", middleware.RequirePermission("items", "read"), handlers.GetItem)
			items.PUT("/:id", middleware.RequirePermission("items", "write"), handlers.UpdateItem)
			items.DELETE("/:id", middleware.RequirePermission("items", "delete"), handlers.DeleteItem)

			// Product Enhancements
			items.POST("/:id/media", middleware.RequirePermission("items", "write"), handlers.UploadItemMedia)
			items.POST("/:id/reviews", middleware.RequirePermission("reviews", "write"), handlers.CreateReview)
			items.POST("/:id/favorite", middleware.RequirePermission("favorites", "write"), handlers.ToggleFavorite)
		}

		// Categories
		categories := api.Group("/categories")
		categories.Use(middleware.AuthMiddleware())
		{
			categories.POST("", middleware.RequirePermission("categories", "write"), handlers.CreateCategory)
			categories.GET("", middleware.RequirePermission("categories", "read"), handlers.GetCategories)
			categories.PUT("/:id", middleware.RequirePermission("categories", "write"), handlers.UpdateCategory)
			categories.DELETE("/:id", middleware.RequirePermission("categories", "delete"), handlers.DeleteCategory)
		}

		// Warehouses
		warehouses := api.Group("/warehouses")
		warehouses.Use(middleware.AuthMiddleware())
		{
			warehouses.POST("", middleware.RequirePermission("warehouses", "write"), handlers.CreateWarehouse)
			warehouses.GET("", middleware.RequirePermission("warehouses", "read"), handlers.GetWarehouses)
			warehouses.PUT("/:id", middleware.RequirePermission("warehouses", "write"), handlers.UpdateWarehouse)
			warehouses.DELETE("/:id", middleware.RequirePermission("warehouses", "delete"), handlers.DeleteWarehouse)
		}

		// Suppliers
		suppliers := api.Group("/suppliers")
		suppliers.Use(middleware.AuthMiddleware())
		{
			suppliers.POST("", middleware.RequirePermission("suppliers", "write"), handlers.CreateSupplier)
			suppliers.GET("", middleware.RequirePermission("suppliers", "read"), handlers.GetSuppliers)
			suppliers.PUT("/:id", middleware.RequirePermission("suppliers", "write"), handlers.UpdateSupplier)
			suppliers.DELETE("/:id", middleware.RequirePermission("suppliers", "delete"), handlers.DeleteSupplier)
		}

		// Discounts
		discounts := api.Group("/discounts")
		discounts.Use(middleware.AuthMiddleware())
		{
			discounts.POST("", middleware.RequirePermission("discounts", "write"), handlers.CreateDiscount)
			discounts.GET("", middleware.RequirePermission("discounts", "read"), handlers.GetDiscounts)
			discounts.PUT("/:id", middleware.RequirePermission("discounts", "write"), handlers.UpdateDiscount)
			discounts.DELETE("/:id", middleware.RequirePermission("discounts", "delete"), handlers.DeleteDiscount)
		}

		// Inventory
		inventory := api.Group("/inventory")
		inventory.Use(middleware.AuthMiddleware())
		{
			inventory.GET("", middleware.RequirePermission("inventory", "read"), handlers.GetInventory)
			inventory.POST("/add", middleware.RequirePermission("inventory", "write"), handlers.AddStock)
			inventory.POST("/transfer", middleware.RequirePermission("inventory", "write"), handlers.TransferStock)
			inventory.PUT("/:id", middleware.RequirePermission("inventory", "write"), handlers.UpdateInventory)
			inventory.DELETE("/:id", middleware.RequirePermission("inventory", "delete"), handlers.DeleteInventory)
		}

		// Purchase Orders
		pos := api.Group("/purchase-orders")
		pos.Use(middleware.AuthMiddleware())
		{
			pos.POST("", middleware.RequirePermission("purchase_orders", "write"), handlers.CreatePurchaseOrder)
			pos.GET("", middleware.RequirePermission("purchase_orders", "read"), handlers.GetPurchaseOrders)
			pos.PUT("/:id/status", middleware.RequirePermission("purchase_orders", "write"), handlers.UpdatePurchaseOrderStatus)
			pos.DELETE("/:id", middleware.RequirePermission("purchase_orders", "delete"), handlers.DeletePurchaseOrder)
		}

		// Sales Orders
		orders := api.Group("/orders")
		orders.Use(middleware.AuthMiddleware())
		{
			orders.POST("", middleware.RequirePermission("orders", "write"), handlers.CreateOrder)
		}

		// Reports & Dashboard
		reports := api.Group("/reports")
		reports.Use(middleware.AuthMiddleware())
		{
			reports.GET("/financial", middleware.RequirePermission("reports", "read"), handlers.GetFinancialReport)
			reports.GET("/sales", middleware.RequirePermission("reports", "read"), handlers.GetSalesReport)
			reports.GET("/dashboard", middleware.RequirePermission("reports", "read"), handlers.GetDashboardSummary)
		}

		// RBAC Management
		rbac := api.Group("/rbac")
		rbac.Use(middleware.AuthMiddleware())
		{
			rbac.POST("/roles", middleware.RequirePermission("roles", "write"), handlers.CreateRole)
			rbac.GET("/roles", middleware.RequirePermission("roles", "read"), handlers.GetRoles)
			rbac.POST("/permissions", middleware.RequirePermission("roles", "write"), handlers.CreatePermission)
			rbac.GET("/permissions", middleware.RequirePermission("roles", "read"), handlers.GetPermissions)
			rbac.POST("/roles/:id/permissions", middleware.RequirePermission("roles", "write"), handlers.AssignPermissionsToRole)
			rbac.POST("/users/:id/role", middleware.RequirePermission("roles", "write"), handlers.AssignRoleToUser)
		}
	}

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
