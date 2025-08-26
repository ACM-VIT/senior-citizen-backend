package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func ExportCSVRouter(app *fiber.App) {
	ExportCSVRouter := app.Group("/export/csv")
	ExportCSVRouter.Get("/", middleware.AuthMiddleware, middleware.AdminMiddleware,
		controllers.ExportMembersCSV)
}
