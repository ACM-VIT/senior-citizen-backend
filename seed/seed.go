package seed

import (
	"acm-portal/database"
	"acm-portal/models"
	"log"
	"time"

	"gorm.io/gorm"
)

// Seed function to populate the database with 4 users, 4 tenures, and 4 memberships
func Seed() {
	// Start a database transaction
	tx := database.DB.Begin()
	defer func() {
		// Handle panic and rollback if necessary
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("Recovered from panic during seeding:", r)
		}
	}()

	// Seed Users with sample data
	users := []models.User{
		{
			Name:              "User 1",
			Email:             "user1@acm.com",
			ProfilePictureURL: "https://example.com/user1.jpg",
			Phone:             "9876543210",
			Gender:            models.Male,
			Bio:               nil,
			TheSentence:       nil,
			LinkedIn:          nil,
			GitHub:            nil,
			Twitter:           nil,
			Instagram:         nil,
			JobTitle:          nil,
			Company:           nil,
			Website:           nil,
		},
		{
			Name:              "User 2",
			Email:             "user2@acm.com",
			ProfilePictureURL: "https://example.com/user2.jpg",
			Phone:             "9876543211",
			Gender:            models.Female,
			Bio:               nil,
			TheSentence:       nil,
			LinkedIn:          nil,
			GitHub:            nil,
			Twitter:           nil,
			Instagram:         nil,
			JobTitle:          nil,
			Company:           nil,
			Website:           nil,
		},
		{
			Name:              "User 3",
			Email:             "user3@acm.com",
			ProfilePictureURL: "https://example.com/user3.jpg",
			Phone:             "9876543212",
			Gender:            models.Male,
			Bio:               nil,
			TheSentence:       nil,
			LinkedIn:          nil,
			GitHub:            nil,
			Twitter:           nil,
			Instagram:         nil,
			JobTitle:          nil,
			Company:           nil,
			Website:           nil,
		},
		{
			Name:              "User 4",
			Email:             "user4@acm.com",
			ProfilePictureURL: "https://example.com/user4.jpg",
			Phone:             "9876543213",
			Gender:            models.Female,
			Bio:               nil,
			TheSentence:       nil,
			LinkedIn:          nil,
			GitHub:            nil,
			Twitter:           nil,
			Instagram:         nil,
			JobTitle:          nil,
			Company:           nil,
			Website:           nil,
		},
	}

	// Insert Users and update the user slice with IDs
	for i, user := range users {
		if err := tx.Create(&user).Error; err != nil {
			rollbackAndLogError(tx, "user", err)
			return
		}
		users[i] = user // Ensure user ID is populated after insertion
		log.Printf("Inserted User: %v", user)
	}

	// Seed Tenures with sample data
	tenures := []models.Tenure{
		{
			IsCurrentTenure: false,
			CurrentYear:     2025,
			StartDate:       time.Now().AddDate(0, -1, 0),
			EndDate:         time.Now().AddDate(0, 0, 0),
		},
		{
			IsCurrentTenure: false,
			CurrentYear:     2025,
			StartDate:       time.Now().AddDate(0, -2, 0),
			EndDate:         time.Now().AddDate(0, -1, 0),
		},
		{
			IsCurrentTenure: false,
			CurrentYear:     2025,
			StartDate:       time.Now().AddDate(0, -3, 0),
			EndDate:         time.Now().AddDate(0, -2, 0),
		},
		{
			IsCurrentTenure: false,
			CurrentYear:     2025,
			StartDate:       time.Now().AddDate(0, -4, 0),
			EndDate:         time.Now().AddDate(0, -3, 0),
		},
	}

	// Insert Tenures and update the tenure slice with IDs
	for i, tenure := range tenures {
		if err := tx.Create(&tenure).Error; err != nil {
			rollbackAndLogError(tx, "tenure", err)
			return
		}
		tenures[i] = tenure // Ensure tenure ID is populated after insertion
		log.Printf("Inserted Tenure: %v", tenure)
	}

	// Seed Memberships with sample data
	memberships := []models.Membership{
		{
			UserID:   users[0].ID,      // Assuming user IDs are UUIDs
			TenureID: tenures[0].ID,    // Assuming tenure IDs are UUIDs
			Role:     models.Community, // Example role
			Position: nil,
		},
		{
			UserID:   users[1].ID,
			TenureID: tenures[1].ID,
			Role:     models.Community,
			Position: nil,
		},
		{
			UserID:   users[2].ID,
			TenureID: tenures[2].ID,
			Role:     models.Community,
			Position: nil,
		},
		{
			UserID:   users[3].ID,
			TenureID: tenures[3].ID,
			Role:     models.Community,
			Position: nil,
		},
	}

	// Insert Memberships
	for _, membership := range memberships {
		if err := tx.Create(&membership).Error; err != nil {
			rollbackAndLogError(tx, "membership", err)
			return
		}
		log.Printf("Inserted Membership: %v", membership)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Println("Error committing transaction:", err)
		tx.Rollback()
		return
	}
	log.Println("Database seeded successfully with 4 users, 4 tenures, and 4 memberships!")
}

// rollbackAndLogError is a helper function to handle rollback and logging errors.
func rollbackAndLogError(tx *gorm.DB, entity string, err error) {
	tx.Rollback()
	log.Printf("Error seeding %s: %v", entity, err)
}
