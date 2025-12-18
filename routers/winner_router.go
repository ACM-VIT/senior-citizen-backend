package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func WinnerRouter(app *fiber.App) {
	winnerGroup := app.Group("/api/v1/winners")
	winnerGroup.Get("/", middleware.AuthMiddleware, controllers.ListWinners)
	winnerGroup.Get("/:id", middleware.AuthMiddleware, controllers.GetWinner)
	winnerGroup.Post("/", middleware.AuthMiddleware, middleware.AdminMiddleware, controllers.CreateWinner)
	winnerGroup.Put("/:id", middleware.AuthMiddleware, middleware.AdminMiddleware, controllers.UpdateWinner)
	winnerGroup.Delete("/:id", middleware.AuthMiddleware, middleware.AdminMiddleware, controllers.DeleteWinner)
}
