package handlers

import (
	"go-rest/internal/database"
	"go-rest/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreatePurchaseOrder(c *gin.Context) {
	var input struct {
		SupplierID  string `json:"supplier_id"`
		WarehouseID string `json:"warehouse_id"`
		Items       []struct {
			ItemID    string  `json:"item_id"`
			Quantity  int     `json:"quantity"`
			UnitPrice float64 `json:"unit_price"`
		} `json:"items"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate total amount
	var totalAmount float64
	var poItems []models.PurchaseOrderItem
	for _, item := range input.Items {
		totalAmount += float64(item.Quantity) * item.UnitPrice
		poItems = append(poItems, models.PurchaseOrderItem{
			ItemID:    uuid.MustParse(item.ItemID),
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})
	}

	po := models.PurchaseOrder{
		SupplierID:  uuid.MustParse(input.SupplierID),
		WarehouseID: uuid.MustParse(input.WarehouseID),
		Status:      "Pending",
		TotalAmount: totalAmount,
		Date:        time.Now(),
		Items:       poItems,
	}

	if err := database.DB.Create(&po).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, po)
}

func UpdatePurchaseOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var po models.PurchaseOrder
	if err := database.DB.Preload("Items").First(&po, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase Order not found"})
		return
	}

	if po.Status == "Received" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PO already received"})
		return
	}

	// If status changes to Received, update inventory
	if input.Status == "Received" {
		err := database.DB.Transaction(func(tx *gorm.DB) error {
			for _, item := range po.Items {
				var inventory models.Inventory
				if err := tx.Where("item_id = ? AND warehouse_id = ?", item.ItemID, po.WarehouseID).First(&inventory).Error; err != nil {
					// Create if not exists
					inventory = models.Inventory{
						ItemID:      item.ItemID,
						WarehouseID: po.WarehouseID,
						Quantity:    item.Quantity,
					}
					if err := tx.Create(&inventory).Error; err != nil {
						return err
					}
				} else {
					inventory.Quantity += item.Quantity
					if err := tx.Save(&inventory).Error; err != nil {
						return err
					}
				}
			}
			return nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	po.Status = input.Status
	if err := database.DB.Save(&po).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, po)
}
