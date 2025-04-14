package routers

import (
	"acm-portal/controllers"

	"github.com/gofiber/fiber/v2"
)

func SearchMembership(app *fiber.App) {
	SearchMembers := app.Group("/search/member")
	SearchMembers.Get("/", controllers.SearchMembership)
}
