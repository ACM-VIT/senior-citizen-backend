package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func UpdateInformationRouter(app *fiber.App) {
	UpdateInformationRouter := app.Group("/update")
	UpdateInformationRouter.Put("/", middleware.AuthMiddleware, controllers.UpdateInformation)
}
