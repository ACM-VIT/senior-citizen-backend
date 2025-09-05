package controllers

import (
	"acm-portal/database"
	"acm-portal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetOrCreateAdmin(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid or missing userID")
	}

	var admin models.Admin
	result := database.DB.Where("user_id = ?", userID).First(&admin)
	if result.Error == nil {
		return c.JSON(fiber.Map{
			"created": false,
			"admin":   admin,
		})
	}

	if result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	admin = models.Admin{UserID: userID}
	if err := database.DB.Create(&admin).Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"created": true,
		"admin":   admin,
	})
}
