package routers

import (
	"acm-portal/handlers"

	"github.com/gofiber/fiber/v2"
)

func ProtectedRouter(app *fiber.App) {
	protectedGroup := app.Group("/protected")
	protectedGroup.Get("/", handlers.RequireAuth, func(c *fiber.Ctx) error {
		claims := c.Locals("claims")
		return c.JSON(fiber.Map{
			"message": "You are authenticated",
			"claims":  claims,
		})
	})
}
