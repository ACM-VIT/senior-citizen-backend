package routers

import (
	"acm-portal/handlers"

	"github.com/gofiber/fiber/v2"
)

func LoginRouter(app *fiber.App) {
	loginGroup := app.Group("/api/v1/auth/login")
	loginGroup.Get("/", handlers.LoginHandler)
}
