package models

type EmailPreferences struct {
	Base
	UserID          string `gorm:"type:uuid;not null;unique" json:"user_id"`
	User            User   `gorm:"foreignKey:UserID" json:"user"`
	Newsletter      bool   `gorm:"default:true" json:"newsletter"`
	Promotions      bool   `gorm:"default:true" json:"promotions"`
	Blogs           bool   `gorm:"default:true" json:"blogs"`
	EventRecordings bool   `gorm:"default:true" json:"event_recordings"`
}
