package routers

import (
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProtectedRouter(app *fiber.App) {
	protectedGroup := app.Group("/api/v1/auth/protected")
	protectedGroup.Get("/", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		claims := c.Locals("claims")
		return c.JSON(fiber.Map{"message": "You are authenticated", "claims": claims})
	})
}
