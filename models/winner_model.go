package models

type Winner struct {
	Base
	ParticipationID string        `gorm:"type:uuid;not null" json:"participation_id"`
	Participation   Participation `gorm:"foreignKey:ParticipationID" json:"participation"`
	EventID         string        `gorm:"type:uuid;not null;uniqueIndex:idx_event_position" json:"event_id"`
	Position        uint          `gorm:"not null;uniqueIndex:idx_event_position" json:"position"`
}
