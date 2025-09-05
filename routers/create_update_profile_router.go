package routers

import (
	"acm-portal/controllers"

	"github.com/gofiber/fiber/v2"
)

func CreateOrUpdateProfileRouter(app *fiber.App) {
	CreateOrUpdateProfileGroup := app.Group("/profile")
	CreateOrUpdateProfileGroup.Post("/", controllers.CreateOrUpdateProfile)
}
