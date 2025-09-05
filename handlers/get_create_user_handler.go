package handlers

import (
	"acm-portal/controllers"

	"github.com/gofiber/fiber/v2"
)

func GetOrCreateUserHandler(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "email is required")
	}

	userID, err := controllers.GetOrCreateUser(email)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"user_id": userID,
		"email":   email,
	})
}