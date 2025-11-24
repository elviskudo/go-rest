package models

type Supplier struct {
	Base
	Name        string `json:"name" form:"name"`
	ContactInfo string `json:"contact_info" form:"contact_info"`
	Address     string `json:"address" form:"address"`
}
