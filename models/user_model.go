package models

type User struct {
	Base
	Name              string `gorm:"type:varchar(100);not null;" json:"name"`
	Email             string `gorm:"type:varchar(100);not null;uniqueIndex" json:"email"`
	ProfilePictureURL string `gorm:"type:text" json:"profile_picture_url"`
	Phone             string `gorm:"uniqueIndex;type:varchar(20);not null" json:"phone"`
	Gender            Gender `gorm:"type:text;not null" json:"gender"`
}

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
	Other  Gender = "other"
)
