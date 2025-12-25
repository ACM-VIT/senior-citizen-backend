package models

import "github.com/lib/pq"

type Profile struct {
	Base
	UserID            string `gorm:"uniqueIndex" json:"user_id"`
	Name              string `gorm:"type:varchar(100);not null;" json:"name"`
	ProfilePictureURL string `gorm:"type:text" json:"profile_picture_url"`
	Phone             string `gorm:"uniqueIndex:idx_users_phone;not null" json:"phone"`
	Gender            Gender `gorm:"type:text;not null" json:"gender"`
	Bio               string `gorm:"type:text" json:"bio,omitempty"`
	TheSentence       string `gorm:"type:text" json:"the_sentence,omitempty"`
	LinkedIn          string `gorm:"type:varchar(255)" json:"linkedin,omitempty"`
	GitHub            string `gorm:"type:varchar(255)" json:"github,omitempty"`
	Twitter           string `gorm:"type:varchar(255)" json:"twitter,omitempty"`
	Instagram         string `gorm:"type:varchar(255)" json:"instagram,omitempty"`
	JobTitle          string `gorm:"type:varchar(100)" json:"job_title,omitempty"`
	Company           string `gorm:"type:varchar(100)" json:"company,omitempty"`
	Website           string `gorm:"type:varchar(255)" json:"website,omitempty"`
	// Columns to store Trigram indexes, used for search
	NameGrams    pq.StringArray `gorm:"type:text[]" json:"-"`
	CompanyGrams pq.StringArray `gorm:"type:text[]" json:"-"`
}

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
	Other  Gender = "other"
)
