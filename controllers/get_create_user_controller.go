package controllers

//a few changes, error handling needs to be better
import (
	"acm-portal/database"
	"acm-portal/models"

	"gorm.io/gorm"
)

func FetchOrCreateUser(email, name, phone, gender string) (string, error) {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error == nil {
		return user.ID, nil
	}

	if result.Error != gorm.ErrRecordNotFound {
		return "", result.Error
	}

	user = models.User{
		Email:  email,
		Name:   name,
		Phone:  phone,
		Gender: models.Gender(gender),
	}
	err := database.DB.Create(&user).Error
	return user.ID, err
}
