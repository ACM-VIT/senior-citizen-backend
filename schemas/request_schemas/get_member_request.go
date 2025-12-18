package requestschemas

import "acm-portal/models"

type GetMembersRequest struct {
	Role     models.Role `query:"role"`
	Position string      `query:"position"`
	Page     int         `query:"page"`
	PageSize int         `query:"page_size"`
}
