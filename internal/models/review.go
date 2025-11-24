package models

import "github.com/google/uuid"

type Review struct {
	Base
	UserID  uuid.UUID `json:"user_id"`
	ItemID  uuid.UUID `json:"item_id"`
	Rating  int       `json:"rating"`
	Comment string    `json:"comment"`
}
