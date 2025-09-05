package routers

import (
	"acm-portal/controllers"

	"github.com/gofiber/fiber/v2"
)

func SearchMembershipRouter(app *fiber.App) {
	SearchMembersGroup := app.Group("/search/member")
	SearchMembersGroup.Post("/", controllers.SearchMembership)
}
