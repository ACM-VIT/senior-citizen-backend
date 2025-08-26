package routers

import (
	"acm-portal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SignupRouter(app *fiber.App) {
	SignupRouter := app.Group("/login")
	SignupRouter.Post("/", handlers.SignupOrLoginHandler)
}
