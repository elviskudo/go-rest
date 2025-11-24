package models

import "github.com/google/uuid"

type User struct {
	Base
	Username string    `gorm:"unique" json:"username"`
	Password string    `json:"-"`
	RoleID   uuid.UUID `json:"role_id"`
	Role     Role      `json:"role"`
}
