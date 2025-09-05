package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func GetCurrentUserRouter(app *fiber.App) {
	GetCurrentUserGroup := app.Group("/current/user")
	GetCurrentUserGroup.Get("/", middleware.AuthMiddleware,
		controllers.GetCurrentUser)
}
