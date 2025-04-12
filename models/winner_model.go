package models

type Winner struct {
	Base
	ParticipationID string        `gorm:"type:uuid;not null" json:"participation_id"`
	Participation   Participation `gorm:"foreignKey:ParticipationID" json:"participation"`
	Position        uint          `gorm:"not null;uniqueIndex:idx_event_position" json:"position"`
}