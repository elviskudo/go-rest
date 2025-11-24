package handlers

import (
	"go-rest/internal/database"
	"go-rest/internal/models"
	"go-rest/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AddStock adds stock to a warehouse
// AddStock godoc
// @Summary      Add stock
// @Description  Add stock to inventory
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        input  body      object  true  "Stock Input"
// @Success      201    {object}  models.Inventory
// @Failure      400    {object}  gin.H
// @Failure      500    {object}  gin.H
// @Security     BearerAuth
// @Router       /inventory/add [post]
func AddStock(c *gin.Context) {
	var input struct {
		ItemID      string `json:"item_id"`
		WarehouseID string `json:"warehouse_id"`
		Quantity    int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var inventory models.Inventory
	if err := database.DB.Where("item_id = ? AND warehouse_id = ?", input.ItemID, input.WarehouseID).First(&inventory).Error; err != nil {
		// Create new inventory record
		inventory = models.Inventory{
			ItemID:      uuid.MustParse(input.ItemID),
			WarehouseID: uuid.MustParse(input.WarehouseID),
			Quantity:    input.Quantity,
		}
		if err := database.DB.Create(&inventory).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Update existing record
		inventory.Quantity += input.Quantity
		if err := database.DB.Save(&inventory).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, inventory)
}

// TransferStock moves stock from one warehouse to another
// TransferStock godoc
// @Summary      Transfer stock
// @Description  Transfer stock between warehouses
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        input  body      object  true  "Transfer Input"
// @Success      200    {object}  gin.H
// @Failure      400    {object}  gin.H
// @Failure      500    {object}  gin.H
// @Security     BearerAuth
// @Router       /inventory/transfer [post]
func TransferStock(c *gin.Context) {
	var input struct {
		ItemID          string `json:"item_id"`
		FromWarehouseID string `json:"from_warehouse_id"`
		ToWarehouseID   string `json:"to_warehouse_id"`
		Quantity        int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Transaction
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var fromInventory models.Inventory
		if err := tx.Where("item_id = ? AND warehouse_id = ?", input.ItemID, input.FromWarehouseID).First(&fromInventory).Error; err != nil {
			return err
		}

		if fromInventory.Quantity < input.Quantity {
			return gorm.ErrInvalidData // Not enough stock
		}

		fromInventory.Quantity -= input.Quantity
		if err := tx.Save(&fromInventory).Error; err != nil {
			return err
		}

		var toInventory models.Inventory
		if err := tx.Where("item_id = ? AND warehouse_id = ?", input.ItemID, input.ToWarehouseID).First(&toInventory).Error; err != nil {
			toInventory = models.Inventory{
				ItemID:      uuid.MustParse(input.ItemID),
				WarehouseID: uuid.MustParse(input.ToWarehouseID),
				Quantity:    input.Quantity,
			}
			if err := tx.Create(&toInventory).Error; err != nil {
				return err
			}
		} else {
			toInventory.Quantity += input.Quantity
			if err := tx.Save(&toInventory).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock transferred successfully"})
}

// GetInventory godoc
// @Summary      List inventory
// @Description  Get inventory items with filters
// @Tags         inventory
// @Produce      json
// @Param        warehouse_id  query     string  false  "Warehouse ID"
// @Param        item_id       query     string  false  "Item ID"
// @Param        page          query     int     false  "Page number"
// @Param        page_size     query     int     false  "Page size"
// @Success      200           {array}   models.Inventory
// @Failure      500           {object}  gin.H
// @Security     BearerAuth
// @Router       /inventory [get]
func GetInventory(c *gin.Context) {
	var inventory []models.Inventory
	query := database.DB.Model(&models.Inventory{})

	// Filter by Warehouse if provided
	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		query = query.Where("warehouse_id = ?", warehouseID)
	}

	// Filter by Item if provided
	if itemID := c.Query("item_id"); itemID != "" {
		query = query.Where("item_id = ?", itemID)
	}

	query = query.Scopes(utils.Paginate(c))

	if err := query.Preload("Item").Preload("Warehouse").Find(&inventory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

// UpdateInventory godoc
// @Summary      Update inventory
// @Description  Update inventory quantity manually (Correction)
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "Inventory ID"
// @Param        input  body      object  true  "Quantity Input"
// @Success      200    {object}  models.Inventory
// @Failure      400    {object}  gin.H
// @Failure      404    {object}  gin.H
// @Failure      500    {object}  gin.H
// @Security     BearerAuth
// @Router       /inventory/{id} [put]
func UpdateInventory(c *gin.Context) {
	id := c.Param("id")
	var inventory models.Inventory
	if err := database.DB.First(&inventory, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory record not found"})
		return
	}

	var input struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inventory.Quantity = input.Quantity

	if err := database.DB.Save(&inventory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

// DeleteInventory godoc
// @Summary      Delete inventory
// @Description  Delete an inventory record
// @Tags         inventory
// @Produce      json
// @Param        id   path      string  true  "Inventory ID"
// @Success      200  {object}  gin.H
// @Failure      404  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Security     BearerAuth
// @Router       /inventory/{id} [delete]
func DeleteInventory(c *gin.Context) {
	id := c.Param("id")
	var inventory models.Inventory
	if err := database.DB.First(&inventory, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory record not found"})
		return
	}

	if err := database.DB.Delete(&inventory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory record deleted successfully"})
}
