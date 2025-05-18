package routers

import (
	"acm-portal/controllers"

	"github.com/gofiber/fiber/v2"
)

func UpsertMembershipRouter(app *fiber.App) {
	Membership := app.Group("/upsert/membership")
	Membership.Post("/", controllers.UpsertMembership)
}
