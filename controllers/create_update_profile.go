package controllers

import (
	"strings"

	"acm-portal/database"
	"acm-portal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateOrUpdateProfile(c *fiber.Ctx) error {
	var body profileUpdateInput

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Token invalid or expired.")
	}

	db := database.DB

	var existing models.Profile
	err := db.Where("user_id = ?", userID).First(&existing).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch profile"})
	}

	if err == gorm.ErrRecordNotFound {
		if body.Name == nil || strings.TrimSpace(*body.Name) == "" ||
			body.Phone == nil || strings.TrimSpace(*body.Phone) == "" ||
			body.Gender == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name, phone, and gender are required"})
		}

		if *body.Gender != models.Male && *body.Gender != models.Female && *body.Gender != models.Other {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid gender"})
		}

		profile := models.Profile{
			UserID:            userID,
			Name:              strings.TrimSpace(*body.Name),
			ProfilePictureURL: "",
			Phone:             strings.TrimSpace(*body.Phone),
			Gender:            *body.Gender,
		}
		if body.ProfilePictureURL != nil {
			profile.ProfilePictureURL = strings.TrimSpace(*body.ProfilePictureURL)
		}
		if body.Bio != nil {
			profile.Bio = strings.TrimSpace(*body.Bio)
		}
		if body.TheSentence != nil {
			profile.TheSentence = strings.TrimSpace(*body.TheSentence)
		}
		if body.LinkedIn != nil {
			profile.LinkedIn = strings.TrimSpace(*body.LinkedIn)
		}
		if body.GitHub != nil {
			profile.GitHub = strings.TrimSpace(*body.GitHub)
		}
		if body.Twitter != nil {
			profile.Twitter = strings.TrimSpace(*body.Twitter)
		}
		if body.Instagram != nil {
			profile.Instagram = strings.TrimSpace(*body.Instagram)
		}
		if body.JobTitle != nil {
			profile.JobTitle = strings.TrimSpace(*body.JobTitle)
		}
		if body.Company != nil {
			profile.Company = strings.TrimSpace(*body.Company)
		}
		if body.Website != nil {
			profile.Website = strings.TrimSpace(*body.Website)
		}

		if err := db.Create(&profile).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(profile)
	}

	updates, err := buildProfileUpdates(body)
	if err != nil {
		return err
	}
	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no valid fields to update"})
	}

	if err := db.Model(&existing).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := db.Where("user_id = ?", userID).First(&existing).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch updated profile"})
	}

	return c.JSON(existing)
}
