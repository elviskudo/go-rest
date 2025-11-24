package models

import "time"

type Discount struct {
	Base
	Name       string    `json:"name"`
	Percentage float64   `json:"percentage"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Active     bool      `json:"active"`
}
