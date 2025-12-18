package requestschemas

type UpdateIndividualInput struct {
	UserID      string `json:"user_id,omitempty" validate:"omitempty,uuid4"`
	Email       string `json:"email" validate:"required,email"`
	Name        string `json:"name,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Gender      string `json:"gender,omitempty"`
	TenureID    string `json:"tenure_id" validate:"required,uuid4"`
	CurrentYear int    `json:"current_year,omitempty"`
	Role        string `json:"role,omitempty"`
	Position    string `json:"position,omitempty"`
}
