package models

import "github.com/google/uuid"

type Item struct {
	Base
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`

	CategoryID uuid.UUID `json:"category_id"`
	SupplierID uuid.UUID `json:"supplier_id"`

	// Quantity removed, moved to Inventory
	ViewerCount   int `json:"viewer_count"`
	FavoriteCount int `json:"favorite_count"`

	Media    []Media   `json:"media"`
	Variants []Variant `json:"variants"`
	Reviews  []Review  `json:"reviews"`
}
