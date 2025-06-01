package models

type Membership struct {
	Base
	UserID   string `gorm:"type:uuid;not null;uniqueIndex:idx_user_tenure"`
	User     User   `gorm:"foreignKey:UserID" json:"user"`
	TenureID string `gorm:"type:uuid;not null;uniqueIndex:idx_user_tenure"`
	Tenure   Tenure `gorm:"foreignKey:TenureID" json:"tenure"`
	Role     Role   `gorm:"type:varchar(20);not null" json:"role"`
	Position string `gorm:"type:varchar(100)" json:"position,omitempty"`
}

type Role string

const (
	Community Role = "community"
	JC        Role = "jc"
	SC        Role = "sc"
	Board     Role = "board"
)