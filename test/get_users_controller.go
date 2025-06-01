package test

import (
	"acm-portal/database"
	"acm-portal/models"

	"github.com/gofiber/fiber/v2"
)

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch users.")
	}
	return c.JSON(users)
}
