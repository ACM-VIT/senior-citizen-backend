package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func ParticipationRouter(app *fiber.App) {
	participationGroup := app.Group("/api/v1/participations")
	participationGroup.Get("/", middleware.AuthMiddleware, controllers.ListParticipations)
	participationGroup.Get("/:id", middleware.AuthMiddleware, controllers.GetParticipation)
	participationGroup.Post("/", middleware.AuthMiddleware, controllers.CreateParticipation)
	participationGroup.Put("/:id", middleware.AuthMiddleware, controllers.UpdateParticipation)
	participationGroup.Delete("/:id", middleware.AuthMiddleware, controllers.DeleteParticipation)
}
