package database

import (
	"acm-portal/models"
	"log"
	"strings"

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
		if isIgnorableMigrationError(err) {
			log.Println("Migration warning:", err)
			return
		}
		log.Fatal("Failed to run migrations: ", err)
	} else {
		log.Println("Migrations completed successfully.")
	}
}

func isIgnorableMigrationError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "constraint \"uni_users_email\"") &&
		strings.Contains(msg, "does not exist")
}
