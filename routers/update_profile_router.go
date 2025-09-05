package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func UpdateProfileRouter(app *fiber.App) {
	UpdateProfileGroup := app.Group("/api/v1/profile/update")
	UpdateProfileGroup.Put("/", middleware.AuthMiddleware, controllers.UpdateProfile)
}
