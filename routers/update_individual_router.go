package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func UpdateIndividualRouter(app *fiber.App) {
	Membership := app.Group("/update/individual")
	Membership.Post("/", middleware.AuthMiddleware, controllers.UpsertIndividual)
}
