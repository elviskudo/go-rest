package models

type User struct {
	Base
	Username string `gorm:"unique" json:"username"`
	Password string `json:"-"`
}
