package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func UpdateProfileRouter(app *fiber.App) {
	UpdateProfileGroup := app.Group("/update/profile")
	UpdateProfileGroup.Put("/", middleware.AuthMiddleware, controllers.UpdateProfile)
}
