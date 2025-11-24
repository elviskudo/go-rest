package models

import "github.com/google/uuid"

type Media struct {
	Base
	ItemID   uuid.UUID `json:"item_id"`
	Type     string    `json:"type"` // "image" or "video"
	URL      string    `json:"url"`
	PublicID string    `json:"public_id"`
}
