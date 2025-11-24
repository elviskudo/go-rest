package models

import "github.com/google/uuid"

type User struct {
	Base
	Username string    `gorm:"unique" json:"username" form:"username"`
	Password string    `json:"-" form:"password"`
	RoleID   uuid.UUID `json:"role_id" form:"role_id"`
	Role     Role      `json:"role"`
}
