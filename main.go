package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"acm-portal/database"
	"acm-portal/routers"
	"acm-portal/test"
)

func initApp() {
	database.ConnectToDB()
	if os.Getenv("SHOULD_MIGRATE") == "true" {
		database.RunMigrations(database.DB)
		// database.SetupPgTrgm(database.DB)
	}
	// test.Test()
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

	test.GetAllUsersRouter(app)
	test.GetAllUsersRouter(app)
	routers.CallbackRouter(app)
	routers.ConsentRouter(app)
	routers.CreateOrUpdateProfileRouter(app)
	routers.ExportCSVRouter(app)
	routers.GetAdminRouter(app)
	routers.CreateAdminRouter(app)
	routers.GetOrCreateUserRouter(app)
	routers.GetCurrentUserRouter(app)
	routers.GetMembersRouter(app)
	routers.ProtectedRouter(app)
	routers.SearchMembershipRouter(app)
	routers.UpdateProfileRouter(app)
	routers.UpdateRolesRouter(app)

}
