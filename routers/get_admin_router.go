package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func GetAdminRouter(app *fiber.App) {
	GetAdminGroup := app.Group("/api/v1/admin/get")
	GetAdminGroup.Get("/", middleware.AuthMiddleware, controllers.GetAdmin)
}
