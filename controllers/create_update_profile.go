package controllers

import (
	"acm-portal/database"
	"acm-portal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateOrUpdateProfile(c *fiber.Ctx) error {
	var body models.Profile

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if body.Gender != "" && body.Gender != models.Male && body.Gender != models.Female && body.Gender != models.Other {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid gender"})
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
		if body.Name == "" || body.Phone == "" || body.Gender == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name, phone, and gender are required"})
		}

		body.UserID = userID
		if err := db.Create(&body).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(body)
	}

	body.UserID = userID
	if err := db.Model(&existing).Updates(body).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := db.Where("user_id = ?", userID).First(&existing).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch updated profile"})
	}

	return c.JSON(existing)
}
