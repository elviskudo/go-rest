package models

type Category struct {
	Base
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
}
