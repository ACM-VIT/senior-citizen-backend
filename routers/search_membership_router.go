package routers

import (
	"acm-portal/controllers"

	"github.com/gofiber/fiber/v2"
)

func SearchMembershipRouter(app *fiber.App) {
	SearchMembers := app.Group("/search/member")
	SearchMembers.Get("/", controllers.SearchMembership)
}
