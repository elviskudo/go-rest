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

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
		return
	}

	var revenue float64
	database.DB.Model(&models.Order{}).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Select("sum(total_amount)").
		Scan(&revenue)

	var cost float64
	database.DB.Model(&models.PurchaseOrder{}).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Select("sum(total_amount)").
		Scan(&cost)

	c.JSON(http.StatusOK, gin.H{
		"revenue":    revenue,
		"cost":       cost,
		"net_profit": revenue - cost,
	})
}

// GetSalesReport godoc
// @Summary      Get sales report
// @Description  Get daily sales aggregation
// @Tags         reports
// @Produce      json
// @Success      200  {array}   object
// @Failure      500  {object}  gin.H
// @Security     BearerAuth
// @Router       /reports/sales [get]
func GetSalesReport(c *gin.Context) {
	type SalesData struct {
		Date       string  `json:"date"`
		TotalSales float64 `json:"total_sales"`
		OrderCount int     `json:"order_count"`
	}

	var sales []SalesData
	database.DB.Model(&models.Order{}).
		Select("date(date) as date, sum(total_amount) as total_sales, count(id) as order_count").
		Group("date(date)").
		Scan(&sales)

	c.JSON(http.StatusOK, sales)
}
