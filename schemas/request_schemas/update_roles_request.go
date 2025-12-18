package requestschemas

type UpdateRolesRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required"`
	Position string `json:"position,omitempty"`
	Tenure   string `json:"tenure" validate:"required"`
}
