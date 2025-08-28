// routers/membership.go
package routers

import (
	"acm-portal/handlers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func UpdateRolesRouter(app *fiber.App) {
	app.Post("/update/roles", middleware.AuthMiddleware, middleware.AdminMiddleware, handlers.UpdateRolesHandler)
}
