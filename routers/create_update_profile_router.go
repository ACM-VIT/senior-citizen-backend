package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func CreateOrUpdateProfileRouter(app *fiber.App) {
	CreateOrUpdateProfileGroup := app.Group("/api/v1/profile/create")
	CreateOrUpdateProfileGroup.Post("/", middleware.AuthMiddleware, controllers.CreateOrUpdateProfile)
}
