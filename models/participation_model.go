package models

type Participation struct {
	Base
	UserID   string  `gorm:"type:uuid;not null" json:"user_id"`
	User     User    `gorm:"foreignKey:UserID" json:"user"`
	EventID  string  `gorm:"type:uuid;not null" json:"event_id"`
	Event    Event   `gorm:"foreignKey:EventID" json:"event"`
	TeamName *string `gorm:"type:varchar(255)" json:"team_name,omitempty"`
	IsWinner bool    `gorm:"default:false" json:"is_winner"`
}
