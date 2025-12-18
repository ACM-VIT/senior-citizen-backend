package routers

import (
	"acm-portal/handlers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func GetOrCreateUserRouter(app *fiber.App) {
	GetOrCreateUserGroup := app.Group("/api/v1/user/create")
	GetOrCreateUserGroup.Get("/", middleware.AuthMiddleware, middleware.AdminMiddleware, handlers.GetOrCreateUserHandler)
}
