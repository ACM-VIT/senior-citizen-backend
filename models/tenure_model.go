package models

import (
	"time"
)

type Tenure struct {
	Base
	IsCurrentTenure bool      `gorm:"default:false" json:"is_current_tenure"`
	CurrentYear     int       `json:"current_year"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
}
