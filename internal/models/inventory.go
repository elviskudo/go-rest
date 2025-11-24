package models

import "github.com/google/uuid"

type Inventory struct {
	Base
	ItemID      uuid.UUID `json:"item_id"`
	WarehouseID uuid.UUID `json:"warehouse_id"`
	Quantity    int       `json:"quantity"`
}
