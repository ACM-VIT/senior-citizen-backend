package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func GetMembersRouter(app *fiber.App) {
	GetMembersRouter := app.Group("/member")
	// GetMembersRouter.Use(middleware.AdminMiddleware())
	GetMembersRouter.Get("/",
		middleware.AuthMiddleware,
		controllers.GetMembers)
}
