package models

type Participation struct {
	Base
	UserID   string  `gorm:"type:uuid;not null;uniqueIndex:idx_user_event" json:"user_id"`
	User     User    `gorm:"foreignKey:UserID" json:"user"`
	EventID  string  `gorm:"type:uuid;not null;uniqueIndex:idx_user_event" json:"event_id"`
	Event    Event   `gorm:"foreignKey:EventID" json:"event"`
	TeamName *string `gorm:"type:varchar(255)" json:"team_name,omitempty"`
}
