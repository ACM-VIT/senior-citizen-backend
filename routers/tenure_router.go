package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func TenureRouter(app *fiber.App) {
	tenureGroup := app.Group("/api/v1/tenures")
	tenureGroup.Get("/", middleware.AuthMiddleware, controllers.ListTenures)
	tenureGroup.Get("/:id", middleware.AuthMiddleware, controllers.GetTenure)
	tenureGroup.Post("/", middleware.AuthMiddleware, middleware.AdminMiddleware, controllers.CreateTenure)
	tenureGroup.Put("/:id", middleware.AuthMiddleware, middleware.AdminMiddleware, controllers.UpdateTenure)
	tenureGroup.Delete("/:id", middleware.AuthMiddleware, middleware.AdminMiddleware, controllers.DeleteTenure)
}
