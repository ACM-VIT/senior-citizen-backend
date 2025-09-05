package routers

import (
	"acm-portal/handlers"

	"github.com/gofiber/fiber/v2"
)

func CallbackRouter(app *fiber.App) {
	callbackGroup := app.Group("/api/v1/auth/callback")
	callbackGroup.Get("/", handlers.CallbackHandler)
}