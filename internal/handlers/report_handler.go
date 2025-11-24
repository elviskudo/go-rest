package handlers

import (
	"go-rest/internal/database"
	"go-rest/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetFinancialReport(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr == "" {
		startDate = time.Now().AddDate(0, 0, -30) // Default last 30 days
	} else {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format (YYYY-MM-DD)"})
			return
		}
	}

	if endDateStr == "" {
		endDate = time.Now()
	} else {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format (YYYY-MM-DD)"})
			return
		}
	}

	// Revenue (Sales)
	var totalRevenue float64
	database.DB.Model(&models.Order{}).
		Where("date BETWEEN ? AND ? AND status = ?", startDate, endDate, "Completed").
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&totalRevenue)

	// Costs (Purchase Orders)
	var totalCost float64
	database.DB.Model(&models.PurchaseOrder{}).
		Where("date BETWEEN ? AND ? AND status != ?", startDate, endDate, "Cancelled").
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&totalCost)

	profit := totalRevenue - totalCost

	c.JSON(http.StatusOK, gin.H{
		"start_date":    startDate.Format("2006-01-02"),
		"end_date":      endDate.Format("2006-01-02"),
		"total_revenue": totalRevenue,
		"total_cost":    totalCost,
		"net_profit":    profit,
	})
}

func GetSalesReport(c *gin.Context) {
	// Aggregate sales by day
	type DailySales struct {
		Date       string  `json:"date"`
		TotalSales float64 `json:"total_sales"`
		OrderCount int     `json:"order_count"`
	}

	var sales []DailySales

	// SQLite specific date function
	database.DB.Model(&models.Order{}).
		Select("strftime('%Y-%m-%d', date) as date, SUM(total_amount) as total_sales, COUNT(id) as order_count").
		Where("status = ?", "Completed").
		Group("date").
		Order("date DESC").
		Scan(&sales)

	c.JSON(http.StatusOK, sales)
}
