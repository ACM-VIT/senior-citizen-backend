package routers

import (
	"acm-portal/controllers"

	"github.com/gofiber/fiber/v2"
)

func ExportCSVRouter(app *fiber.App) {
	ExportCSVRouter := app.Group("/export/csv")
	ExportCSVRouter.Get("/", controllers.ExportMembersCSV)
}