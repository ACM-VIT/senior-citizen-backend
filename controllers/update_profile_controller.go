package controllers

import (
	"strings"

	"acm-portal/database"
	"acm-portal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type profileUpdateInput struct {
	Name              *string        `json:"name"`
	ProfilePictureURL *string        `json:"profile_picture_url"`
	Phone             *string        `json:"phone"`
	Gender            *models.Gender `json:"gender"`
	Bio               *string        `json:"bio"`
	TheSentence       *string        `json:"the_sentence"`
	LinkedIn          *string        `json:"linkedin"`
	GitHub            *string        `json:"github"`
	Twitter           *string        `json:"twitter"`
	Instagram         *string        `json:"instagram"`
	JobTitle          *string        `json:"job_title"`
	Company           *string        `json:"company"`
	Website           *string        `json:"website"`
}

func UpdateProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Token invalid or expired.")
	}

	var profile models.Profile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Profile not found.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch profile. Please try again later.")
	}

	var updateData profileUpdateInput
	if err := c.BodyParser(&updateData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body provided. Please check your request format.")
	}

	updates := map[string]interface{}{}
	if updateData.Name != nil {
		trimmed := strings.TrimSpace(*updateData.Name)
		if trimmed == "" {
			return fiber.NewError(fiber.StatusBadRequest, "name cannot be empty")
		}
		updates["name"] = trimmed
	}
	if updateData.ProfilePictureURL != nil {
		updates["profile_picture_url"] = strings.TrimSpace(*updateData.ProfilePictureURL)
	}
	if updateData.Phone != nil {
		trimmed := strings.TrimSpace(*updateData.Phone)
		if trimmed == "" {
			return fiber.NewError(fiber.StatusBadRequest, "phone cannot be empty")
		}
		updates["phone"] = trimmed
	}
	if updateData.Gender != nil {
		if *updateData.Gender != models.Male && *updateData.Gender != models.Female && *updateData.Gender != models.Other {
			return fiber.NewError(fiber.StatusBadRequest, "invalid gender")
		}
		updates["gender"] = *updateData.Gender
	}
	if updateData.Bio != nil {
		updates["bio"] = strings.TrimSpace(*updateData.Bio)
	}
	if updateData.TheSentence != nil {
		updates["the_sentence"] = strings.TrimSpace(*updateData.TheSentence)
	}
	if updateData.LinkedIn != nil {
		updates["linked_in"] = strings.TrimSpace(*updateData.LinkedIn)
	}
	if updateData.GitHub != nil {
		updates["git_hub"] = strings.TrimSpace(*updateData.GitHub)
	}
	if updateData.Twitter != nil {
		updates["twitter"] = strings.TrimSpace(*updateData.Twitter)
	}
	if updateData.Instagram != nil {
		updates["instagram"] = strings.TrimSpace(*updateData.Instagram)
	}
	if updateData.JobTitle != nil {
		updates["job_title"] = strings.TrimSpace(*updateData.JobTitle)
	}
	if updateData.Company != nil {
		updates["company"] = strings.TrimSpace(*updateData.Company)
	}
	if updateData.Website != nil {
		updates["website"] = strings.TrimSpace(*updateData.Website)
	}

	if len(updates) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "no valid fields to update")
	}

	if err := database.DB.Model(&profile).Updates(updates).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update profile.")
	}

	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch updated profile.")
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Information updated successfully", "data": profile})
}
