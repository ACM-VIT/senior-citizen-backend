package routers

import (
	"acm-portal/handlers"

	"github.com/gofiber/fiber/v2"
)

func ConsentRouter(app *fiber.App) {
	consentGroup := app.Group("/consent")
	consentGroup.Get("/", handlers.ConsentHandler)
}
