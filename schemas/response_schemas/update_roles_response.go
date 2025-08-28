package responseschemas

type UpdateRolesResponse struct {
	Email    string `json:"email"`
	Role     string `json:"role"`
	Position string `json:"position,omitempty"`
	Tenure   string `json:"tenure"`
	Status   string `json:"status"`
}