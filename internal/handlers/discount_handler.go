package handlers

import (
	"go-rest/internal/database"
	"go-rest/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateDiscount godoc
// @Summary      Create a discount
// @Description  Create a new discount
// @Tags         discounts
// @Accept       json
// @Produce      json
// @Param        discount  body      models.Discount  true  "Discount JSON"
// @Success      201       {object}  models.Discount
// @Failure      400       {object}  gin.H
// @Failure      500       {object}  gin.H
// @Security     BearerAuth
// @Router       /discounts [post]
func CreateDiscount(c *gin.Context) {
	var discount models.Discount
	if err := c.ShouldBindJSON(&discount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&discount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, discount)
}

// GetDiscounts godoc
// @Summary      List discounts
// @Description  Get all discounts
// @Tags         discounts
// @Produce      json
// @Success      200  {array}   models.Discount
// @Failure      500  {object}  gin.H
// @Security     BearerAuth
// @Router       /discounts [get]
func GetDiscounts(c *gin.Context) {
	var discounts []models.Discount
	if err := database.DB.Find(&discounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, discounts)
}

// UpdateDiscount godoc
// @Summary      Update a discount
// @Description  Update a discount by ID
// @Tags         discounts
// @Accept       json
// @Produce      json
// @Param        id        path      string           true  "Discount ID"
// @Param        discount  body      models.Discount  true  "Discount JSON"
// @Success      200       {object}  models.Discount
// @Failure      400       {object}  gin.H
// @Failure      404       {object}  gin.H
// @Failure      500       {object}  gin.H
// @Security     BearerAuth
// @Router       /discounts/{id} [put]
func UpdateDiscount(c *gin.Context) {
	id := c.Param("id")
	var discount models.Discount
	if err := database.DB.First(&discount, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Discount not found"})
		return
	}

	var input models.Discount
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	discount.Name = input.Name
	discount.Percentage = input.Percentage
	discount.StartDate = input.StartDate
	discount.EndDate = input.EndDate
	discount.Active = input.Active

	if err := database.DB.Save(&discount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, discount)
}

// DeleteDiscount godoc
// @Summary      Delete a discount
// @Description  Delete a discount by ID
// @Tags         discounts
// @Produce      json
// @Param        id   path      string  true  "Discount ID"
// @Success      200  {object}  gin.H
// @Failure      404  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Security     BearerAuth
// @Router       /discounts/{id} [delete]
func DeleteDiscount(c *gin.Context) {
	id := c.Param("id")
	var discount models.Discount
	if err := database.DB.First(&discount, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Discount not found"})
		return
	}

	if err := database.DB.Delete(&discount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Discount deleted successfully"})
}
