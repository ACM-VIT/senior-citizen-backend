package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"acm-portal/database"
	"acm-portal/routers"
	"acm-portal/seed"
)

func initApp() {
	database.ConnectToDB()
	if os.Getenv("SHOULD_MIGRATE") == "true" {
		database.RunMigrations(database.DB)
	}
	// seed.Seed()
}

func main() {
	initApp()

	app := fiber.New()

	setupRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber + GORM!")
	})

	log.Fatal(app.Listen(":3000"))
}

func setupRoutes(app *fiber.App) {

	routers.GetMembersRouter(app)
	routers.ExportCSVRouter(app)
	seed.GetAllUsersRouter(app)
	routers.UpdateInformationRouter(app)
	routers.SearchMembershipRouter(app)
}
