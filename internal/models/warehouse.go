package models

type Warehouse struct {
	Base
	Name     string `json:"name" form:"name"`
	Location string `json:"location" form:"location"`
	Capacity int    `json:"capacity" form:"capacity"`
}
