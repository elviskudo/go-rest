package models

import "github.com/google/uuid"

type Favorite struct {
	Base
	UserID uuid.UUID `json:"user_id"`
	ItemID uuid.UUID `json:"item_id"`
}
