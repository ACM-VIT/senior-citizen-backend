package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func EmailPreferencesRouter(app *fiber.App) {
	prefsGroup := app.Group("/api/v1/email-preferences")
	prefsGroup.Get("/", middleware.AuthMiddleware, controllers.GetEmailPreferences)
	prefsGroup.Put("/", middleware.AuthMiddleware, controllers.UpdateEmailPreferences)
}
