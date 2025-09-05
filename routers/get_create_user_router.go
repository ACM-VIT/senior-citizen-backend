package routers

import (
	"acm-portal/handlers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func GetOrCreateUserRouter(app *fiber.App) {
	GetOrCreateUserGroup := app.Group("/user")
	GetOrCreateUserGroup.Get("/",middleware.AuthMiddleware, handlers.GetOrCreateUserHandler)
}
