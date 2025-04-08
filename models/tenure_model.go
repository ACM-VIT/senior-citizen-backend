package models

import (
	"time"
)

type Tenure struct {
	Base
	BoardHandoverOn    *time.Time `gorm:"type:timestamp;unique" json:"board_handover_on"`
	CommunityCreatedAt *time.Time `gorm:"type:timestamp" json:"community_created_at"`
	CoreInductedAt     *time.Time `gorm:"type:timestamp" json:"core_inducted_at"`
	CoreExitedAt       *time.Time `gorm:"type:timestamp" json:"core_exited_at,omitempty"`
	IsCurrentTenure    bool       `gorm:"default:false" json:"is_current_tenure"`
}