package requestschemas

import "acm-portal/models"

type UpsertMembershipInput struct {
	UserID   string `json:"user_id" validate:"required,uuid4""`
	TenureID string `json:"tenure_id" validate:"required,uuid4""`
	Role        models.Role  `json:"role,omitempty"`         
    Position    string       `json:"position,omitempty"` 
}