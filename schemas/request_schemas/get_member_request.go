package requestschemas

import "acm-portal/models"

type GetMembersRequest struct {
    Gender            *models.Gender `json:"gender,omitempty"`                         
    JobTitle          *string        `json:"job_title,omitempty"`           
    Company           *string        `json:"company,omitempty"`                   
    CreatedAfter      *string        `json:"created_after,omitempty"`       
    UpdatedAfter      *string        `json:"updated_after,omitempty"`       
    Page              *int           `json:"page,omitempty"`                
    PageSize          *int           `json:"page_size,omitempty"`           
}