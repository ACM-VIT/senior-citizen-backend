package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SearchMembershipRouter(app *fiber.App) {
	SearchMembersGroup := app.Group("/api/v1/search/member")
	SearchMembersGroup.Post("/", middleware.AuthMiddleware, controllers.SearchMembership)
}
