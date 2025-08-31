package controllers

import (
	"acm-portal/database"
	"acm-portal/models"

	"github.com/gofiber/fiber/v2"
)

func GetCurrentUser(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Token invalid or expired.")
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}
