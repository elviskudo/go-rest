package handlers

import (
	"go-rest/internal/database"
	"go-rest/internal/models"
	"go-rest/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateSupplier(c *gin.Context) {
	var supplier models.Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, supplier)
}

func GetSuppliers(c *gin.Context) {
	var suppliers []models.Supplier
	query := database.DB.Model(&models.Supplier{})

	query = query.Scopes(utils.Search(c, []string{"name", "contact_info", "address"}))
	query = query.Scopes(utils.Sort(c, map[string]bool{"name": true}))
	query = query.Scopes(utils.Paginate(c))

	if err := query.Find(&suppliers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, suppliers)
}
