package handlers

import (
	"go-rest/internal/database"
	"go-rest/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateCategory godoc
// @Summary      Create a category
// @Description  Create a new product category
// @Tags         categories
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        name         formData  string  true  "Category Name"
// @Param        description  formData  string  false "Category Description"
// @Success      201          {object}  models.Category
// @Failure      400          {object}  gin.H
// @Failure      500          {object}  gin.H
// @Security     BearerAuth
// @Router       /categories [post]
func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBind(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetCategories godoc
// @Summary      List categories
// @Description  Get all product categories
// @Tags         categories
// @Produce      json
// @Success      200  {array}   models.Category
// @Failure      500  {object}  gin.H
// @Security     BearerAuth
// @Router       /categories [get]
func GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := database.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// UpdateCategory godoc
// @Summary      Update a category
// @Description  Update a product category by ID
// @Tags         categories
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        id           path      string  true  "Category ID"
// @Param        name         formData  string  true  "Category Name"
// @Param        description  formData  string  false "Category Description"
// @Success      200          {object}  models.Category
// @Failure      400          {object}  gin.H
// @Failure      404          {object}  gin.H
// @Failure      500          {object}  gin.H
// @Security     BearerAuth
// @Router       /categories/{id} [put]
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := database.DB.First(&category, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var input models.Category
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.Name = input.Name
	category.Description = input.Description

	if err := database.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory godoc
// @Summary      Delete a category
// @Description  Delete a product category by ID
// @Tags         categories
// @Produce      json
// @Param        id   path      string  true  "Category ID"
// @Success      200  {object}  gin.H
// @Failure      404  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Security     BearerAuth
// @Router       /categories/{id} [delete]
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := database.DB.First(&category, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if err := database.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
