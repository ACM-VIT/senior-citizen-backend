package test

import (
	"github.com/gofiber/fiber/v2"
)

func GetAllUsersRouter(app *fiber.App) {
	GetAllUsersRouter := app.Group("/members")
	GetAllUsersRouter.Get("/", GetAllUsers)
}