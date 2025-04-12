package responseschemas

import "acm-portal/models"

type GetMembersResponse struct {
    ID                string          `json:"id"`
    Name              string          `json:"name"`
    Email             string          `json:"email"`
    Phone             string          `json:"phone"`
    Gender            models.Gender  `json:"gender"`
    ProfilePictureURL string          `json:"profile_picture_url"`
}

type PaginatedMembersResponse struct {
    Data       []GetMembersResponse `json:"data"`
    Total      int64                `json:"total"`
    Page       int                  `json:"page"`
    PageSize   int                  `json:"page_size"`
    TotalPages int                  `json:"total_pages"`
}
