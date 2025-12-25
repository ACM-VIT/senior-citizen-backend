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
	}

	if err := createFuzzySearchIndexes(db); err != nil {
		log.Fatal("Failed to create GIN indexes: ", err)
	}

	log.Println("Migrations completed successfully.")
}

func createFuzzySearchIndexes(db *gorm.DB) error {
	return db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_profiles_name_grams
		ON profiles USING GIN (name_grams);

		CREATE INDEX IF NOT EXISTS idx_profiles_company_grams
		ON profiles USING GIN (company_grams);
	`).Error
}

func isIgnorableMigrationError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "constraint \"uni_users_email\"") &&
		strings.Contains(msg, "does not exist")
}
