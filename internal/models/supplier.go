package models

type Supplier struct {
	Base
	Name        string `json:"name"`
	ContactInfo string `json:"contact_info"`
	Address     string `json:"address"`
}
