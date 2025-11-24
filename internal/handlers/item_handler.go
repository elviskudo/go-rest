package handlers

import (
	"go-rest/internal/database"
	"go-rest/internal/models"
	"go-rest/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateItem godoc
// @Summary      Create a new item
// @Description  Create a new inventory item
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        item  body      models.Item  true  "Item JSON"
// @Success      201   {object}  models.Item
// @Failure      400   {object}  gin.H
// @Failure      500   {object}  gin.H
// @Security     BearerAuth
// @Router       /items [post]
func CreateItem(c *gin.Context) {
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// GetItems godoc
// @Summary      List items
// @Description  Get all inventory items with pagination, search, and sort
// @Tags         items
// @Produce      json
// @Param        page       query     int     false  "Page number"
// @Param        page_size  query     int     false  "Page size"
// @Param        search     query     string  false  "Search term"
// @Param        sort       query     string  false  "Sort field"
// @Param        order      query     string  false  "Sort order (asc/desc)"
// @Success      200  {array}   models.Item
// @Failure      500  {object}  gin.H
// @Security     BearerAuth
// @Router       /items [get]
func GetItems(c *gin.Context) {
	var items []models.Item
	query := database.DB.Model(&models.Item{})

	// Search
	query = query.Scopes(utils.Search(c, []string{"name", "description"}))

	// Sort
	query = query.Scopes(utils.Sort(c, map[string]bool{"name": true, "price": true, "created_at": true}))

	// Pagination
	query = query.Scopes(utils.Paginate(c))

	if err := query.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetItem godoc
// @Summary      Get an item
// @Description  Get an inventory item by ID
// @Tags         items
// @Produce      json
// @Param        id   path      string  true  "Item ID"
// @Success      200  {object}  models.Item
// @Failure      404  {object}  gin.H
// @Security     BearerAuth
// @Router       /items/{id} [get]
func GetItem(c *gin.Context) {
	id := c.Param("id")
	var item models.Item

	if err := database.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// UpdateItem godoc
// @Summary      Update an item
// @Description  Update an inventory item by ID
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        id    path      string       true  "Item ID"
// @Param        item  body      models.Item  true  "Item JSON"
// @Success      200   {object}  models.Item
// @Failure      400   {object}  gin.H
// @Failure      404   {object}  gin.H
// @Security     BearerAuth
// @Router       /items/{id} [put]
func UpdateItem(c *gin.Context) {
	id := c.Param("id")
	var item models.Item

	if err := database.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	var input models.Item
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	item.Name = input.Name
	item.Description = input.Description
	item.Price = input.Price
	// Quantity is now managed via Inventory

	if err := database.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteItem godoc
// @Summary      Delete an item
// @Description  Delete an inventory item by ID
// @Tags         items
// @Produce      json
// @Param        id   path      string  true  "Item ID"
// @Success      200  {object}  gin.H
// @Failure      404  {object}  gin.H
// @Security     BearerAuth
// @Router       /items/{id} [delete]
func DeleteItem(c *gin.Context) {
	id := c.Param("id")
	var item models.Item

	if err := database.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	if err := database.DB.Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
