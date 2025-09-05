package controllers

import (
	"acm-portal/database"
	"acm-portal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAdmin(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid or missing userID")
	}

	var admin models.Admin
	result := database.DB.Where("user_id = ?", userID).First(&admin)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User is not an admin"})
		}
		return result.Error
	}

	return c.JSON(admin)
}
