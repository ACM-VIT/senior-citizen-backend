package middleware

import (
	"acm-portal/helpers"

	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or missing userID")
	}

	isAdmin, err := helpers.IsAdmin(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User is not an admin"})
	}

	return c.Next()
}
