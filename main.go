package main

import (
	"acm-portal/database"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectToDB()
	if os.Getenv("SHOULD_MIGRATE") == "true" {
		database.RunMigrations(database.DB)
	}
	app := fiber.New()
	// routers.GetMembersRouter(app)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber + GORM!")
	})
	log.Fatal(app.Listen(":3000"))
}
