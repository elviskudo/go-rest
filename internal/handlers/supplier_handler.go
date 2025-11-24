package handlers

import (
	"go-rest/internal/database"
	"go-rest/internal/models"
	"go-rest/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateSupplier godoc
// @Summary      Create a supplier
// @Description  Create a new supplier
// @Tags         suppliers
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        name          formData  string  true  "Supplier Name"
// @Param        contact_info  formData  string  true  "Contact Info"
// @Param        address       formData  string  true  "Address"
// @Success      201           {object}  models.Supplier
// @Failure      400           {object}  gin.H
// @Failure      500           {object}  gin.H
// @Security     BearerAuth
// @Router       /suppliers [post]
func CreateSupplier(c *gin.Context) {
	var supplier models.Supplier
	if err := c.ShouldBind(&supplier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, supplier)
}

// GetSuppliers godoc
// @Summary      List suppliers
// @Description  Get all suppliers with pagination, search, and sort
// @Tags         suppliers
// @Produce      json
// @Param        page       query     int     false  "Page number"
// @Param        page_size  query     int     false  "Page size"
// @Param        search     query     string  false  "Search term"
// @Param        sort       query     string  false  "Sort field"
// @Param        order      query     string  false  "Sort order (asc/desc)"
// @Success      200  {array}   models.Supplier
// @Failure      500  {object}  gin.H
// @Security     BearerAuth
// @Router       /suppliers [get]
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

// UpdateSupplier godoc
// @Summary      Update a supplier
// @Description  Update a supplier by ID
// @Tags         suppliers
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        id            path      string  true  "Supplier ID"
// @Param        name          formData  string  true  "Supplier Name"
// @Param        contact_info  formData  string  true  "Contact Info"
// @Param        address       formData  string  true  "Address"
// @Success      200           {object}  models.Supplier
// @Failure      400           {object}  gin.H
// @Failure      404           {object}  gin.H
// @Failure      500           {object}  gin.H
// @Security     BearerAuth
// @Router       /suppliers/{id} [put]
func UpdateSupplier(c *gin.Context) {
	id := c.Param("id")
	var supplier models.Supplier
	if err := database.DB.First(&supplier, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		return
	}

	var input models.Supplier
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	supplier.Name = input.Name
	supplier.ContactInfo = input.ContactInfo
	supplier.Address = input.Address

	if err := database.DB.Save(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, supplier)
}

// DeleteSupplier godoc
// @Summary      Delete a supplier
// @Description  Delete a supplier by ID
// @Tags         suppliers
// @Produce      json
// @Param        id   path      string  true  "Supplier ID"
// @Success      200  {object}  gin.H
// @Failure      404  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Security     BearerAuth
// @Router       /suppliers/{id} [delete]
func DeleteSupplier(c *gin.Context) {
	id := c.Param("id")
	var supplier models.Supplier
	if err := database.DB.First(&supplier, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		return
	}

	if err := database.DB.Delete(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Supplier deleted successfully"})
}
