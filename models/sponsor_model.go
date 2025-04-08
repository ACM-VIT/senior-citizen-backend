package models

type Sponsor struct {
	Base
	Name    string  `gorm:"type:varchar(255);not null" json:"name"`
	LogoURL *string `gorm:"type:text" json:"logo_url,omitempty"`
	Website *string `gorm:"type:text" json:"website,omitempty"`
	Contact *string `gorm:"type:varchar(50)" json:"contact,omitempty"`
	//Sponsor money? Not sure
	Events []Event `gorm:"many2many:event_sponsors;" json:"-"`
}