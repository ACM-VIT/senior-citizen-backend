package responseschemas

import "acm-portal/models"

type GetMembersResponse struct {
    ID        string          `json:"id"`          
    UserID    string          `json:"user_id"`     
    User      models.User     `json:"user"`        
    TenureID  string          `json:"tenure_id"`  
    Tenure    models.Tenure  `json:"tenure"`      
    Role      models.Role     `json:"role"`        
    Position  string         `json:"position,omitempty"`
}

type PaginatedMembersResponse struct {
    Data       []GetMembersResponse `json:"data"`        
    Total      int64                   `json:"total"`       
    Page       int                     `json:"page"`        
    PageSize   int                     `json:"page_size"`   
    TotalPages int                     `json:"total_pages"` 
}