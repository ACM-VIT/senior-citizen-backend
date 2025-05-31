package controllers

import (
	"acm-portal/database"
	"acm-portal/models"
	"fmt"
)

func fetchTenureByYear(year int) (string, error) {
	var tenure models.Tenure
	err := database.DB.Where("current_year = ?", year).First(&tenure).Error
	if err != nil {
		return "", fmt.Errorf("tenure not found for year: %d", year)
	}
	return tenure.ID, nil
}
