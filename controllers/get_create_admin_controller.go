package controllers

import (
	"acm-portal/database"
	"acm-portal/models"
)

func FetchOrCreateAdmin(userID string) (bool, error) {
	var admin models.Admin
	if err := database.DB.Where("user_id = ?", userID).First(&admin).Error; err != nil {
		return false, nil
	}
	return true, nil
}
