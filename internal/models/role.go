package models

import (
	"github.com/google/uuid"
)

type Role struct {
	Base
	Name        string       `json:"name" gorm:"unique"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
}

type Permission struct {
	Base
	Resource string `json:"resource"` // e.g., items, users
	Action   string `json:"action"`   // e.g., read, write, delete
}

// Join table for Role <-> Permission
type RolePermission struct {
	RoleID       uuid.UUID `gorm:"primaryKey"`
	PermissionID uuid.UUID `gorm:"primaryKey"`
}
