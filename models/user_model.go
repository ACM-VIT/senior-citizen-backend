package models

type User struct {
	Base
	Email   string  `gorm:"type:varchar(100);not null;uniqueIndex:idx_users_email" json:"email"`
	Profile Profile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}