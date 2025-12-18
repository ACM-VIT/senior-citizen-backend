package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"acm-portal/database"
	"acm-portal/internal/config"
	"acm-portal/routers"
)

func initApp(cfg config.Config) {
	database.ConnectToDB(cfg)
	if cfg.ShouldMigrate {
		database.RunMigrations(database.DB)
	}
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	initApp(cfg)

	app := fiber.New()

	setupRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber + GORM!")
	})

	log.Fatal(app.Listen(":" + cfg.Port))
}

func setupRoutes(app *fiber.App) {

	routers.ConsentRouter(app)
	routers.CreateOrUpdateProfileRouter(app)
	routers.ExportCSVRouter(app)
	routers.GetAdminRouter(app)
	routers.CreateAdminRouter(app)
	routers.GetOrCreateUserRouter(app)
	routers.GetCurrentUserRouter(app)
	routers.GetMembersRouter(app)
	routers.TenureRouter(app)
	routers.EventRouter(app)
	routers.SponsorRouter(app)
	routers.ParticipationRouter(app)
	routers.WinnerRouter(app)
	routers.EmailPreferencesRouter(app)
	routers.ProtectedRouter(app)
	routers.SearchMembershipRouter(app)
	routers.UpdateProfileRouter(app)
	routers.UpdateRolesRouter(app)

}
