package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SponsorRouter(app *fiber.App) {
	sponsorGroup := app.Group("/api/v1/sponsors")
	sponsorGroup.Get("/", middleware.AuthMiddleware, controllers.ListSponsors)
	sponsorGroup.Get("/:id", middleware.AuthMiddleware, controllers.GetSponsor)
	sponsorGroup.Post("/", middleware.AuthMiddleware, middleware.AdminMiddleware, controllers.CreateSponsor)
	sponsorGroup.Put("/:id", middleware.AuthMiddleware, middleware.AdminMiddleware, controllers.UpdateSponsor)
	sponsorGroup.Delete("/:id", middleware.AuthMiddleware, middleware.AdminMiddleware, controllers.DeleteSponsor)
}
