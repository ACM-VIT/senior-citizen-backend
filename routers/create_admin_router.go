package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func CreateAdminRouter(app *fiber.App) {
	CreateAdminGroup := app.Group("/api/v1/admin/create")
	CreateAdminGroup.Post("/", middleware.AuthMiddleware, middleware.AdminMiddleware, controllers.CreateAdmin)
}
