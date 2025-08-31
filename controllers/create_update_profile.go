package controllers

import (
	"acm-portal/database"
	"acm-portal/models"

	"github.com/gofiber/fiber/v2"
)

func CreateProfile(c *fiber.Ctx) error {
	var body models.Profile

	// Parse request body
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// UserID required
	if body.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}

	db := database.DB

	// Check if profile already exists
	var existing models.Profile
	if err := db.Where("user_id = ?", body.UserID).First(&existing).Error; err == nil {
		// Profile already exists
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Profile already exists for this user_id",
		})
	}

	// Create new profile
	if err := db.Create(&body).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(body)
}
