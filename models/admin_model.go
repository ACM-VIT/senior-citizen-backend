package models

type Admin struct {
	Base
	UserID string `gorm:"unique" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID" json:"user"`
}
