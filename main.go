package main

import (
	"acm-portal/database"
	"acm-portal/routers"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectToDB()
	if os.Getenv("SHOULD_MIGRATE") == "true" {
		database.RunMigrations(database.DB)
	}
	//if os.Getenv("SHOULD_SEED") == "true" {
	//	seed.Seed()
	//}
	app := fiber.New()
	routers.GetMembersRouter(app)
	routers.ExportCSVRouter(app)
	routers.SearchMembership(app)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber + GORM!")
	})
	log.Fatal(app.Listen(":3000"))
}
