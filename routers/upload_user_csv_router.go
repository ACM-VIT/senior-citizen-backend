package routers

import (
	"acm-portal/controllers"

	"github.com/gofiber/fiber/v2"
)

func UploadUserCSVRouter(app *fiber.App) {
	userCSVController := controllers.NewUserCSVController()
	app.Post("/users/upload-csv", userCSVController.UploadUsersCSV)
}
