package routers

import (
	"acm-portal/controllers"

	"github.com/gofiber/fiber/v2"
)

func UpdateIndividualRouter(app *fiber.App) {
	Membership := app.Group("/update/individual")
	Membership.Post("/", controllers.UpsertIndividual)
}
