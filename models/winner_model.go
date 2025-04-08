package models

type Winner struct {
	Base
	EventID string  `gorm:"type:uuid;not null;uniqueIndex:idx_event_position"`
	Event   Event   `gorm:"foreignKey:EventID" json:"-"`
	IsTeam  bool    `gorm:"not null;default:false" json:"is_team"`
	UserID  *string `gorm:"type:uuid" json:"user_id,omitempty"`
	User    *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	// TeamID   *string `gorm:"type:uuid" json:"team_id,omitempty"`
	// Team     *Team   `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Position uint `gorm:"not null;uniqueIndex:idx_event_position"`
}
