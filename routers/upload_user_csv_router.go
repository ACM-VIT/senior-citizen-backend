package routers

import (
	"acm-portal/controllers"
	"acm-portal/middleware"

	"github.com/gofiber/fiber/v2"
)

func UploadUserCSVRouter(app *fiber.App) {
	userCSVController := controllers.NewUserCSVController()
	app.Post("/users/upload-csv", middleware.AuthMiddleware, middleware.AdminMiddleware, userCSVController.UploadUsersCSV)
}
