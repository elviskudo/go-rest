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

func CreateOrder(c *gin.Context) {
	var input struct {
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

	// Calculate total amount and prepare items
	var totalAmount float64
	var orderItems []models.OrderItem

	// Transaction to create order and decrease inventory
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range input.Items {
			totalAmount += float64(item.Quantity) * item.UnitPrice
			orderItems = append(orderItems, models.OrderItem{
				ItemID:    uuid.MustParse(item.ItemID),
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
			})

			// Decrease inventory (simplified: assume taking from first available warehouse or specific logic needed)
			// For POS, we usually know the warehouse (store location).
			// Let's assume a default warehouse or pass it in input.
			// For simplicity here, we'll just find ANY inventory and decrease it (FIFO logic could apply but let's keep it simple)
			// Actually, let's just pick the first inventory record with enough stock.

			var inventory models.Inventory
			if err := tx.Where("item_id = ? AND quantity >= ?", item.ItemID, item.Quantity).First(&inventory).Error; err != nil {
				return gorm.ErrInvalidData // Not enough stock
			}

			inventory.Quantity -= item.Quantity
			if err := tx.Save(&inventory).Error; err != nil {
				return err
			}
		}

		order := models.Order{
			UserID:        userID.(uuid.UUID),
			Status:        "Completed",
			PaymentMethod: input.PaymentMethod,
			TotalAmount:   totalAmount,
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
