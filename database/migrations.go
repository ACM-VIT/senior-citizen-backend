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
