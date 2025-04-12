package routers

import (
	"acm-portal/controllers"

	"github.com/gofiber/fiber/v2"
)

func GetMembersRouter(app *fiber.App) {
	GetMembersRouter := app.Group("/admin/members")
	// GetMembersRouter.Use(middleware.AdminMiddleware()) 
	GetMembersRouter.Get("/", controllers.GetMembers)
}