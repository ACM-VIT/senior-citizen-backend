// routers/membership.go
package routers

import (
	"acm-portal/handlers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func UpdateRolesRouter(app *fiber.App) {
	UpdateRolesGroup := app.Group("/update/roles")
	UpdateRolesGroup.Post("/", middleware.AuthMiddleware, middleware.AdminMiddleware, handlers.UpdateRolesHandler)	
}
