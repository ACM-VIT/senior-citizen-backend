package database

import (
	"acm-portal/models"
	"log"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Admin{},
		&models.Base{},
		&models.User{},
		&models.Tenure{},
		&models.Membership{},
		&models.Event{},
		&models.Participation{},
		&models.Sponsor{},
		&models.Winner{},
		&models.EmailPreferences{},
	)

	if err != nil {
		log.Fatal("Failed to run migrations: ", err)
	} else {
		log.Println("Migrations completed successfully.")
	}
}

func SetupPgTrgm(db *gorm.DB) {
	db.Exec(`CREATE EXTENSION IF NOT EXISTS pg_trgm`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_name_trgm ON users USING gin (name gin_trgm_ops)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_memberships_role_trgm ON memberships USING gin (role gin_trgm_ops)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_memberships_position_trgm ON memberships USING gin (position gin_trgm_ops)`)
}
