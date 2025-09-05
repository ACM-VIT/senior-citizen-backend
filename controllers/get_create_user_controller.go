package controllers

import (
	"acm-portal/database"
	"acm-portal/models"

	"gorm.io/gorm"
)

func GetOrCreateUser(email string) (string, error) {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error == nil {
		return user.ID, nil
	}

	if result.Error != gorm.ErrRecordNotFound {
		return "", result.Error
	}

	user = models.User{
		Email: email,
	}
	err := database.DB.Create(&user).Error
	return user.ID, err
}
