package requestschemas

import "acm-portal/models"

type GetMembersRequest struct {
    Role        *models.Role  `json:"role,omitempty"`         
    Position    *string       `json:"position,omitempty"`     
    Page        *int          `json:"page,omitempty"`        
    PageSize    *int          `json:"page_size,omitempty"`    
}