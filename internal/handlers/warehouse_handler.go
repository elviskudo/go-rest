package handlers

import (
	"go-rest/internal/database"
	"go-rest/internal/models"
	"go-rest/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateWarehouse(c *gin.Context) {
	var warehouse models.Warehouse
	if err := c.ShouldBindJSON(&warehouse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&warehouse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, warehouse)
}

func GetWarehouses(c *gin.Context) {
	var warehouses []models.Warehouse
	query := database.DB.Model(&models.Warehouse{})

	query = query.Scopes(utils.Search(c, []string{"name", "location"}))
	query = query.Scopes(utils.Sort(c, map[string]bool{"name": true, "capacity": true}))
	query = query.Scopes(utils.Paginate(c))

	if err := query.Find(&warehouses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, warehouses)
}
