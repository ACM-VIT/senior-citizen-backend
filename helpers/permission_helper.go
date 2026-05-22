package helpers

import (
	"acm-portal/database"
	"acm-portal/models"

	"gorm.io/gorm"
)

type EffectiveAccess struct {
	Role        models.Role
	Permissions []models.PermissionCode
}

func GetCurrentRole(userID string) (models.Role, error) {
	var membership models.Membership
	err := database.DB.Preload("Tenure").Where("user_id = ?", userID).First(&membership).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Community, nil
		}
	}
	return models.Role(membership.Tenure.ID), nil
}
