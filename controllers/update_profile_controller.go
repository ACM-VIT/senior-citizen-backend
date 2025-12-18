package controllers

import (
	"acm-portal/database"
	"acm-portal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

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

	updates, err := buildProfileUpdates(updateData)
	if err != nil {
		return err
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
