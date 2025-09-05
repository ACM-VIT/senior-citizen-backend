package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func GetOrCreateAdminRouter(app *fiber.App) {
	GetOrCreateAdminGroup := app.Group("/admin")
	GetOrCreateAdminGroup.Get("/",middleware.AuthMiddleware, controllers.GetOrCreateAdmin)
}
