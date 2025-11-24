package models

import "github.com/google/uuid"

type Variant struct {
	Base
	ItemID  uuid.UUID `json:"item_id"`
	Name    string    `json:"name"` // e.g., "Size", "Color"
	Options []Option  `json:"options"`
}

type Option struct {
	Base
	VariantID uuid.UUID `json:"variant_id"`
	Name      string    `json:"name"` // e.g., "Small", "Red"
}
