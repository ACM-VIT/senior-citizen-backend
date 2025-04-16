package routers

import (
	"acm-portal/controllers"

	"github.com/gofiber/fiber/v2"
)

func GetMembership(app *fiber.App) {
	GetMembership := app.Group("/check/memberr")
	GetMembership.Get("/", controllers.GetCurrentUser)
}
