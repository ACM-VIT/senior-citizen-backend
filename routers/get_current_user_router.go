package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func GetCurrentUserRouter(app *fiber.App) {
	GetCurrentUserGroup := app.Group("/api/v1/user/get")
	GetCurrentUserGroup.Get("/", middleware.AuthMiddleware, controllers.GetCurrentUser)
}
