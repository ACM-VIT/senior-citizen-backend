package routers

import (
	"acm-portal/handlers"

	"github.com/gofiber/fiber/v2"
)

func CallbackRouter(app *fiber.App) {
	app.Get("/callback", handlers.CallbackHandler)
}
