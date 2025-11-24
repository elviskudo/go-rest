package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	Base
	UserID        uuid.UUID   `json:"user_id"`
	TotalAmount   float64     `json:"total_amount"`
	Status        string      `json:"status"` // Completed, Refunded
	PaymentMethod string      `json:"payment_method"`
	Date          time.Time   `json:"date"`
	Items         []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	Base
	OrderID   uuid.UUID `json:"order_id"`
	ItemID    uuid.UUID `json:"item_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
}
