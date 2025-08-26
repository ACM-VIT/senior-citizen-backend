package controllers

import (
	"acm-portal/database"
	"acm-portal/helpers"
	"acm-portal/models"
	"fmt"

	"gorm.io/gorm"
)

// Upsert membership assigns roles and tenures to users
func UpdateRoles(data []struct {
	Email    string `json:"email"`
	Role     string `json:"role"`
	Position string `json:"position,omitempty"`
	Tenure   string `json:"tenure"`
}) error {
	for _, d := range data {
		// check if user exists
		var user models.User
		if err := database.DB.Where("email = ?", d.Email).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("user not found: %s", d.Email)
			}
			return err
		}

		// fetch or create tenure
		tenure, err := helpers.FetchTenure(d.Tenure)
		if err != nil {
			return err
		}

		// upsert membership
		var membership models.Membership
		err = database.DB.
			Where("user_id = ? AND tenure_id = ?", user.ID, tenure.ID).
			First(&membership).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// create new membership
				membership = models.Membership{
					UserID:   user.ID,
					TenureID: tenure.ID,
					Role:     models.Role(d.Role),
					Position: d.Position,
				}
				if err := database.DB.Create(&membership).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			// update existing membership
			membership.Role = models.Role(d.Role)
			membership.Position = d.Position
			if err := database.DB.Save(&membership).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
