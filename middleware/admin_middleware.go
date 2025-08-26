package middleware

import "github.com/gofiber/fiber/v2"

func AdminMiddleware(c *fiber.Ctx) error {
	isAdmin, ok := c.Locals("admin").(bool)
	if !ok || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin privileges required",
		})
	}
	return c.Next()
}
