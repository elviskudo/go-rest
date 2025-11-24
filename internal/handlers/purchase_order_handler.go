package handlers

import (
	"go-rest/internal/database"
	"go-rest/internal/models"
	"go-rest/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreatePurchaseOrder godoc
// @Summary      Create a purchase order
// @Description  Create a new purchase order
// @Tags         purchase_orders
// @Accept       json
// @Produce      json
// @Param        input  body      object  true  "Purchase Order Input"
// @Success      201    {object}  models.PurchaseOrder
// @Failure      400    {object}  gin.H
// @Failure      500    {object}  gin.H
// @Security     BearerAuth
// @Router       /purchase-orders [post]
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

// UpdatePurchaseOrderStatus godoc
// @Summary      Update purchase order status
// @Description  Update status (e.g., Pending -> Received). Updates inventory if Received.
// @Tags         purchase_orders
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "Purchase Order ID"
// @Param        input  body      object  true  "Status Input"
// @Success      200    {object}  models.PurchaseOrder
// @Failure      400    {object}  gin.H
// @Failure      404    {object}  gin.H
// @Failure      500    {object}  gin.H
// @Security     BearerAuth
// @Router       /purchase-orders/{id}/status [put]
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

	// If status changes to Received, update inventory
	if input.Status == "Received" && po.Status != "Received" {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory"})
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

// GetPurchaseOrders godoc
// @Summary      List purchase orders
// @Description  Get all purchase orders with pagination, search, and sort
// @Tags         purchase_orders
// @Produce      json
// @Param        page       query     int     false  "Page number"
// @Param        page_size  query     int     false  "Page size"
// @Param        search     query     string  false  "Search term"
// @Param        sort       query     string  false  "Sort field"
// @Param        order      query     string  false  "Sort order (asc/desc)"
// @Success      200  {array}   models.PurchaseOrder
// @Failure      500  {object}  gin.H
// @Security     BearerAuth
// @Router       /purchase-orders [get]
func GetPurchaseOrders(c *gin.Context) {
	var pos []models.PurchaseOrder
	query := database.DB.Model(&models.PurchaseOrder{})

	query = query.Scopes(utils.Search(c, []string{"status"})) // Basic search by status
	query = query.Scopes(utils.Sort(c, map[string]bool{"date": true, "total_amount": true}))
	query = query.Scopes(utils.Paginate(c))

	if err := query.Preload("Items").Find(&pos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pos)
}

// DeletePurchaseOrder godoc
// @Summary      Delete a purchase order
// @Description  Delete a purchase order by ID
// @Tags         purchase_orders
// @Produce      json
// @Param        id   path      string  true  "Purchase Order ID"
// @Success      200  {object}  gin.H
// @Failure      404  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Security     BearerAuth
// @Router       /purchase-orders/{id} [delete]
func DeletePurchaseOrder(c *gin.Context) {
	id := c.Param("id")
	var po models.PurchaseOrder
	if err := database.DB.First(&po, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase Order not found"})
		return
	}

	if err := database.DB.Delete(&po).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase Order deleted successfully"})
}
