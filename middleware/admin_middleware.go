package middleware

import (
	"acm-portal/database"
	"acm-portal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AdminMiddleware(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or missing userID")
	}

	var admin models.Admin
	err := database.DB.Where("user_id = ?", userID).First(&admin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User is not an admin"})
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Next()
}
