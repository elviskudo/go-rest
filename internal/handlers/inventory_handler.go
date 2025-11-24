package handlers

import (
	"go-rest/internal/database"
	"go-rest/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AddStock adds stock to a warehouse
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

	itemID := uuid.MustParse(input.ItemID)
	warehouseID := uuid.MustParse(input.WarehouseID)

	var inventory models.Inventory
	if err := database.DB.Where("item_id = ? AND warehouse_id = ?", itemID, warehouseID).First(&inventory).Error; err != nil {
		// Create new inventory record
		inventory = models.Inventory{
			ItemID:      itemID,
			WarehouseID: warehouseID,
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

	itemID := uuid.MustParse(input.ItemID)
	fromID := uuid.MustParse(input.FromWarehouseID)
	toID := uuid.MustParse(input.ToWarehouseID)

	// Transaction
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var source models.Inventory
		if err := tx.Where("item_id = ? AND warehouse_id = ?", itemID, fromID).First(&source).Error; err != nil {
			return err
		}

		if source.Quantity < input.Quantity {
			return gorm.ErrInvalidData // Not enough stock
		}

		source.Quantity -= input.Quantity
		if err := tx.Save(&source).Error; err != nil {
			return err
		}

		var dest models.Inventory
		if err := tx.Where("item_id = ? AND warehouse_id = ?", itemID, toID).First(&dest).Error; err != nil {
			dest = models.Inventory{
				ItemID:      itemID,
				WarehouseID: toID,
				Quantity:    input.Quantity,
			}
			if err := tx.Create(&dest).Error; err != nil {
				return err
			}
		} else {
			dest.Quantity += input.Quantity
			if err := tx.Save(&dest).Error; err != nil {
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
