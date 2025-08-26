package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func GetMembership(app *fiber.App) {
	GetMembership := app.Group("/check/memberr")
	GetMembership.Get("/", middleware.AuthMiddleware,
		controllers.GetCurrentUser)
}
