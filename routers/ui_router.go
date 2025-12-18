package routers

import (
	"net/http"

	"acm-portal/handlers"

	"github.com/gofiber/fiber/v2"
)

func UIRouter(app *fiber.App) {
	app.Get("/ui/config.js", handlers.UIConfigHandler)
	app.Static("/ui/assets", "./ui/assets")

	app.Get("/ui", func(c *fiber.Ctx) error {
		return c.Redirect("/ui/login", http.StatusFound)
	})
	app.Get("/ui/*", func(c *fiber.Ctx) error {
		return c.SendFile("./ui/index.html")
	})
}
