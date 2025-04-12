package requestschemas

import "acm-portal/models"

type GetMembersRequest struct {
    Gender            *models.Gender `json:"gender,omitempty"`               
    Name              *string        `json:"name,omitempty"`                 
    Email             *string        `json:"email,omitempty"`                
    Phone             *string        `json:"phone,omitempty"`                
    Bio               *string        `json:"bio,omitempty"`                  
    LinkedIn          *string        `json:"linkedin,omitempty"`             
    GitHub            *string        `json:"github,omitempty"`               
    Twitter           *string        `json:"twitter,omitempty"`              
    Instagram         *string        `json:"instagram,omitempty"`           
    JobTitle          *string        `json:"job_title,omitempty"`           
    Company           *string        `json:"company,omitempty"`             
    Website           *string        `json:"website,omitempty"`             
    ProfilePictureURL *string        `json:"profile_picture_url,omitempty"`
    TheSentence       *string        `json:"the_sentence,omitempty"`        
    CreatedAfter      *string        `json:"created_after,omitempty"`       
    UpdatedAfter      *string        `json:"updated_after,omitempty"`       
    Page              *int           `json:"page,omitempty"`                
    PageSize          *int           `json:"page_size,omitempty"`           
}