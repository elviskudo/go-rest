package handlers

import (
	"errors"
	"go-rest/internal/database"
	"go-rest/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateOrder(c *gin.Context) {
	var input struct {
		WarehouseID   string `json:"warehouse_id"`
		PaymentMethod string `json:"payment_method"`
		Items         []struct {
			ItemID    string  `json:"item_id"`
			Quantity  int     `json:"quantity"`
			UnitPrice float64 `json:"unit_price"`
		} `json:"items"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	// Transaction to create order and decrease inventory
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var totalAmount float64
		var orderItems []models.OrderItem

		for _, item := range input.Items {
			// Check inventory
			var inventory models.Inventory
			if err := tx.Where("item_id = ? AND warehouse_id = ?", item.ItemID, input.WarehouseID).First(&inventory).Error; err != nil {
				return errors.New("item not found in warehouse")
			}

			if inventory.Quantity < item.Quantity {
				return errors.New("insufficient stock for item: " + item.ItemID)
			}

			// Deduct stock
			inventory.Quantity -= item.Quantity
			if err := tx.Save(&inventory).Error; err != nil {
				return err
			}

			totalAmount += float64(item.Quantity) * item.UnitPrice
			orderItems = append(orderItems, models.OrderItem{
				ItemID:    uuid.MustParse(item.ItemID),
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
			})
		}

		// Create Order
		order := models.Order{
			UserID:        userID.(uuid.UUID),
			WarehouseID:   uuid.MustParse(input.WarehouseID),
			TotalAmount:   totalAmount,
			Status:        "Completed",
			PaymentMethod: input.PaymentMethod,
			Date:          time.Now(),
			Items:         orderItems,
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully"})
}
