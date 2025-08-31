package routers

import (
	"acm-portal/controllers"

	"github.com/gofiber/fiber/v2"
)

func CreateOrUpdateProfile(app *fiber.App) {
	CreateOrUpdateProfileRouter := app.Group("/create/profile")
	CreateOrUpdateProfileRouter.Post("/", controllers.CreateProfile)
}
