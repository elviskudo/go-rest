package handlers

import (
	"go-rest/internal/database"
	"go-rest/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDashboardSummary(c *gin.Context) {
	var itemCount int64
	database.DB.Model(&models.Item{}).Count(&itemCount)

	var warehouseCount int64
	database.DB.Model(&models.Warehouse{}).Count(&warehouseCount)

	var userCount int64
	database.DB.Model(&models.User{}).Count(&userCount)

	var supplierCount int64
	database.DB.Model(&models.Supplier{}).Count(&supplierCount)

	// Low stock items (e.g., < 10)
	var lowStockCount int64
	database.DB.Model(&models.Inventory{}).Where("quantity < ?", 10).Count(&lowStockCount)

	c.JSON(http.StatusOK, gin.H{
		"total_items":      itemCount,
		"total_warehouses": warehouseCount,
		"total_users":      userCount,
		"total_suppliers":  supplierCount,
		"low_stock_items":  lowStockCount,
	})
}
