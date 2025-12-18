package controllers

import (
	"acm-portal/database"
	"acm-portal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type emailPreferencesUpdateInput struct {
	Newsletter      *bool `json:"newsletter"`
	Promotions      *bool `json:"promotions"`
	Blogs           *bool `json:"blogs"`
	EventRecordings *bool `json:"event_recordings"`
}

func GetEmailPreferences(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Token invalid or expired.")
	}

	var prefs models.EmailPreferences
	err := database.DB.Where("user_id = ?", userID).First(&prefs).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch email preferences.")
		}

		prefs = models.EmailPreferences{UserID: userID}
		if err := database.DB.Create(&prefs).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to create email preferences.")
		}
	}

	return c.JSON(prefs)
}

func UpdateEmailPreferences(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Token invalid or expired.")
	}

	var body emailPreferencesUpdateInput
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body provided.")
	}

	var prefs models.EmailPreferences
	if err := database.DB.Where("user_id = ?", userID).First(&prefs).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			prefs = models.EmailPreferences{UserID: userID}
			if err := database.DB.Create(&prefs).Error; err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Failed to create email preferences.")
			}
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch email preferences.")
		}
	}

	updates := map[string]interface{}{}
	if body.Newsletter != nil {
		updates["newsletter"] = *body.Newsletter
	}
	if body.Promotions != nil {
		updates["promotions"] = *body.Promotions
	}
	if body.Blogs != nil {
		updates["blogs"] = *body.Blogs
	}
	if body.EventRecordings != nil {
		updates["event_recordings"] = *body.EventRecordings
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&prefs).Updates(updates).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to update email preferences.")
		}
	}

	if err := database.DB.Where("user_id = ?", userID).First(&prefs).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch updated email preferences.")
	}

	return c.JSON(prefs)
}
