package models

import (
	"time"
)

type Event struct {
	Base
	Name              string     `gorm:"type:varchar(255);not null" json:"name"`
	Description       string     `gorm:"type:text" json:"description,omitempty"`
	TenureID          string     `gorm:"type:uuid;not null" json:"tenure_id"`
	Tenure            Tenure     `gorm:"foreignKey:TenureID" json:"tenure"`
	Venue             string     `gorm:"type:varchar(255)" json:"venue"`
	EventType         EventType  `gorm:"type:varchar(20)" json:"event_type"`
	Date              time.Time  `gorm:"type:timestamp" json:"date"`
	StartTime         *time.Time `gorm:"type:timestamp" json:"start_time,omitempty"`
	EndTime           *time.Time `gorm:"type:timestamp" json:"end_time,omitempty"`
	NumberOfAttendees uint       `json:"number_of_attendees"`
	NumberOfWinners   uint       `json:"number_of_winners"`
	Sponsors          []Sponsor  `gorm:"many2many:event_sponsors;" json:"sponsors,omitempty"`
	Winners           []Winner   `gorm:"foreignKey:ParticipationID" json:"winners,omitempty"`
}

type EventType string

const (
	Online  EventType = "online"
	Offline EventType = "offline"
	Hybrid  EventType = "hybrid"
)
