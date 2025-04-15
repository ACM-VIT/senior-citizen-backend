package routers

import (
	"acm-portal/controllers"

	"github.com/gofiber/fiber/v2"
)

func UpdateInformationRouter(app *fiber.App) {
	UpdateInformationRouter := app.Group("/update")
	UpdateInformationRouter.Put("/:id", controllers.UpdateInformation)
}