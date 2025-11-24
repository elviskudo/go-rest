package handlers

import (
	"go-rest/internal/database"
	"go-rest/internal/models"
	"go-rest/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateWarehouse godoc
// @Summary      Create a warehouse
// @Description  Create a new warehouse
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param        warehouse  body      models.Warehouse  true  "Warehouse JSON"
// @Success      201        {object}  models.Warehouse
// @Failure      400        {object}  gin.H
// @Failure      500        {object}  gin.H
// @Security     BearerAuth
// @Router       /warehouses [post]
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

// GetWarehouses godoc
// @Summary      List warehouses
// @Description  Get all warehouses with pagination, search, and sort
// @Tags         warehouses
// @Produce      json
// @Param        page       query     int     false  "Page number"
// @Param        page_size  query     int     false  "Page size"
// @Param        search     query     string  false  "Search term"
// @Param        sort       query     string  false  "Sort field"
// @Param        order      query     string  false  "Sort order (asc/desc)"
// @Success      200  {array}   models.Warehouse
// @Failure      500  {object}  gin.H
// @Security     BearerAuth
// @Router       /warehouses [get]
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

// UpdateWarehouse godoc
// @Summary      Update a warehouse
// @Description  Update a warehouse by ID
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param        id         path      string            true  "Warehouse ID"
// @Param        warehouse  body      models.Warehouse  true  "Warehouse JSON"
// @Success      200        {object}  models.Warehouse
// @Failure      400        {object}  gin.H
// @Failure      404        {object}  gin.H
// @Failure      500        {object}  gin.H
// @Security     BearerAuth
// @Router       /warehouses/{id} [put]
func UpdateWarehouse(c *gin.Context) {
	id := c.Param("id")
	var warehouse models.Warehouse
	if err := database.DB.First(&warehouse, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Warehouse not found"})
		return
	}

	var input models.Warehouse
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	warehouse.Name = input.Name
	warehouse.Location = input.Location
	warehouse.Capacity = input.Capacity

	if err := database.DB.Save(&warehouse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, warehouse)
}

// DeleteWarehouse godoc
// @Summary      Delete a warehouse
// @Description  Delete a warehouse by ID
// @Tags         warehouses
// @Produce      json
// @Param        id   path      string  true  "Warehouse ID"
// @Success      200  {object}  gin.H
// @Failure      404  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Security     BearerAuth
// @Router       /warehouses/{id} [delete]
func DeleteWarehouse(c *gin.Context) {
	id := c.Param("id")
	var warehouse models.Warehouse
	if err := database.DB.First(&warehouse, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Warehouse not found"})
		return
	}

	if err := database.DB.Delete(&warehouse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Warehouse deleted successfully"})
}
