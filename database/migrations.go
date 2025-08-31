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
		&models.Profile{},
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

	db.Exec(`CREATE INDEX event_name_trgm_idx ON events USING gin (name gin_trgm_ops)`)

	db.Exec(`CREATE INDEX user_name_trgm_idx ON users USING gin (name gin_trgm_ops)`)

	db.Exec(`CREATE INDEX profile_job_title_trgm_idx ON profiles USING gin (job_title gin_trgm_ops)`)
}
