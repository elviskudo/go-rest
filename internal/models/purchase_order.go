package models

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseOrder struct {
	Base
	SupplierID  uuid.UUID           `json:"supplier_id"`
	WarehouseID uuid.UUID           `json:"warehouse_id"`
	Status      string              `json:"status"` // Pending, Received, Cancelled
	TotalAmount float64             `json:"total_amount"`
	Date        time.Time           `json:"date"`
	Items       []PurchaseOrderItem `json:"items" gorm:"foreignKey:PurchaseOrderID"`
}

type PurchaseOrderItem struct {
	Base
	PurchaseOrderID uuid.UUID `json:"purchase_order_id"`
	ItemID          uuid.UUID `json:"item_id"`
	Quantity        int       `json:"quantity"`
	UnitPrice       float64   `json:"unit_price"`
}
