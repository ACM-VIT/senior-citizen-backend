package routers

import (
	"acm-portal/handlers"

	"github.com/gofiber/fiber/v2"
)

func CallbackRouter(app *fiber.App) {
	callbackGroup := app.Group("/callback")
	callbackGroup.Get("/", handlers.CallbackHandler)
}