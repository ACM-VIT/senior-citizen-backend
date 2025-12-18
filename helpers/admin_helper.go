package helpers

import (
	"acm-portal/database"
	"acm-portal/models"

	"gorm.io/gorm"
)

func IsAdmin(userID string) (bool, error) {
	var admin models.Admin
	err := database.DB.Where("user_id = ?", userID).First(&admin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
