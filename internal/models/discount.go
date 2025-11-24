package models

import "time"

type Discount struct {
	Base
	Name       string    `json:"name" form:"name"`
	Percentage float64   `json:"percentage" form:"percentage"`
	StartDate  time.Time `json:"start_date" form:"start_date" time_format:"2006-01-02T15:04:05Z07:00"`
	EndDate    time.Time `json:"end_date" form:"end_date" time_format:"2006-01-02T15:04:05Z07:00"`
	Active     bool      `json:"active" form:"active"`
}
