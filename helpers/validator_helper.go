package helpers

import (
	"acm-portal/database"
	"acm-portal/models"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

func ValidateField(fieldType, value string) error {
	if value == "" {
		if fieldType == "position" {
			return nil
		}
	}

	validValues := map[string]map[string]bool{
		"role": {
			"community": true,
			"jc":        true,
			"sc":        true,
			"board":     true,
		},
		"position": {
			"chairperson":           true,
			"vice-chairperson":      true,
			"secretary":             true,
			"cosecretary":           true,
			"technical-director":    true,
			"projects-lead":         true,
			"research-lead":         true,
			"contest-lead":          true,
			"design-lead":           true,
			"creative-lead":         true,
			"acmw-chairperson":      true,
			"acmw-vice-chairperson": true,
			"acmw-secretary":        true,
		},
	}

	if allowed, exists := validValues[fieldType]; exists {
		if value != "" && !allowed[value] {
			return fmt.Errorf("invalid %s: %s", fieldType, value)
		}
	} else {
		return fmt.Errorf("unknown field type: %s", fieldType)
	}

	return nil
}

func ValidateRole(role, position string) (string, string, error) {
	if err := ValidateField("role", role); err != nil {
		return "", "", err
	}
	if err := ValidateField("position", position); err != nil {
		return "", "", err
	}
	return role, position, nil
}

func FetchTenure(tenureStr string) (*models.Tenure, error) {
	if tenureStr == "" {
		return nil, fmt.Errorf("tenure cannot be empty")
	}

	year, err := strconv.Atoi(tenureStr)
	if err != nil {
		return nil, fmt.Errorf("invalid tenure format: %s", tenureStr)
	}

	var tenure models.Tenure
	err = database.DB.Where("current_year = ?", year).First(&tenure).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			tenure = models.Tenure{
				CurrentYear: year,
			}
			if createErr := database.DB.Create(&tenure).Error; createErr != nil {
				return nil, createErr
			}
			return &tenure, nil
		}
		return nil, err
	}
	return &tenure, nil
}
