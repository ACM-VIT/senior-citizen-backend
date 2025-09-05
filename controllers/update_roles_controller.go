package controllers

import (
	"acm-portal/database"
	"acm-portal/helpers"
	"acm-portal/models"
	requestschemas "acm-portal/schemas/request_schemas"
	responseschemas "acm-portal/schemas/response_schemas"
	"fmt"

	"gorm.io/gorm"
)

func UpdateRoles(data []requestschemas.UpdateRolesRequest) ([]responseschemas.UpdateRolesResponse, error) {
	var results []responseschemas.UpdateRolesResponse

	for _, d := range data {
		role, position, err := helpers.ValidateRole(d.Role, d.Position)
		if err != nil {
			return nil, err
		}

		var user models.User
		if err := database.DB.Where("email = ?", d.Email).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("user not found: %s", d.Email)
			}
			return nil, err
		}

		tenure, err := helpers.FetchTenure(d.Tenure)
		if err != nil {
			return nil, err
		}

		var membership models.Membership
		err = database.DB.
			Where("user_id = ? AND tenure_id = ?", user.ID, tenure.ID).
			First(&membership).Error

		status := "updated"
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				membership = models.Membership{
					UserID:   user.ID,
					TenureID: tenure.ID,
					Role:     models.Role(role),
					Position: position,
				}
				if err := database.DB.Create(&membership).Error; err != nil {
					return nil, err
				}
				status = "created"
			} else {
				return nil, err
			}
		} else {
			membership.Role = models.Role(role)
			membership.Position = position
			if err := database.DB.Save(&membership).Error; err != nil {
				return nil, err
			}
		}

		results = append(results, responseschemas.UpdateRolesResponse{
			Email:    d.Email,
			Role:     role,
			Position: position,
			Tenure:   d.Tenure,
			Status:   status,
		})
	}

	return results, nil
}
