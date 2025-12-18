package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func EventRouter(app *fiber.App) {
	eventGroup := app.Group("/api/v1/events")
	eventGroup.Get("/", middleware.AuthMiddleware, controllers.ListEvents)
	eventGroup.Get("/:id", middleware.AuthMiddleware, controllers.GetEvent)
	eventGroup.Post("/", middleware.AuthMiddleware, middleware.AdminMiddleware, controllers.CreateEvent)
	eventGroup.Put("/:id", middleware.AuthMiddleware, middleware.AdminMiddleware, controllers.UpdateEvent)
	eventGroup.Delete("/:id", middleware.AuthMiddleware, middleware.AdminMiddleware, controllers.DeleteEvent)
}
