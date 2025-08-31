package test

import (
	"acm-portal/database"
	"acm-portal/models"
	"log"
	"time"

	"gorm.io/gorm"
)

func Test() {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("Recovered from panic during seeding:", r)
		}
	}()

	users := []models.User{
		{Email: "user1@acm.com"},
		{Email: "user2@acm.com"},
		{Email: "user3@acm.com"},
		{Email: "user4@acm.com"},
		{Email: "user5@acm.com"},
	}

	for i := range users {
		if err := tx.Create(&users[i]).Error; err != nil {
			rollbackAndLogError(tx, "user", err)
			return
		}
		log.Printf("Inserted User: %v", users[i])
	}

	profiles := []models.Profile{
		{
			UserID:           users[0].ID,
			Name:             "User 1",
			ProfilePictureURL: "https://example.com/user1.jpg",
			Phone:            "9876543210",
			Gender:           models.Male,
			Bio:              "This is User 1's bio.",
			TheSentence:      "User 1's life motto.",
			LinkedIn:         "https://linkedin.com/in/user1",
			GitHub:           "https://github.com/user1",
			Twitter:          "https://twitter.com/user1",
			Instagram:        "https://instagram.com/user1",
			JobTitle:         "Software Engineer",
			Company:          "Tech Corp",
			Website:          "https://user1.com",
		},
		{
			UserID:           users[1].ID,
			Name:             "User 2",
			ProfilePictureURL: "https://example.com/user2.jpg",
			Phone:            "9876543211",
			Gender:           models.Female,
			Bio:              "This is User 2's bio.",
			TheSentence:      "User 2's life motto.",
			LinkedIn:         "https://linkedin.com/in/user2",
			GitHub:           "https://github.com/user2",
			Twitter:          "https://twitter.com/user2",
			Instagram:        "https://instagram.com/user2",
			JobTitle:         "Data Scientist",
			Company:          "Data Solutions Inc.",
			Website:          "https://user2.com",
		},
		{
			UserID:           users[2].ID,
			Name:             "User 3",
			ProfilePictureURL: "https://example.com/user3.jpg",
			Phone:            "9876543212",
			Gender:           models.Male,
			Bio:              "This is User 3's bio.",
			TheSentence:      "User 3's life motto.",
			LinkedIn:         "https://linkedin.com/in/user3",
			GitHub:           "https://github.com/user3",
			Twitter:          "https://twitter.com/user3",
			Instagram:        "https://instagram.com/user3",
			JobTitle:         "UI/UX Designer",
			Company:          "Design Studios",
			Website:          "https://user3.com",
		},
		{
			UserID:           users[3].ID,
			Name:             "User 4",
			ProfilePictureURL: "https://example.com/user4.jpg",
			Phone:            "9876543213",
			Gender:           models.Female,
			TheSentence:      "User 4's life motto.",
			LinkedIn:         "https://linkedin.com/in/user4",
			GitHub:           "https://github.com/user4",
			Twitter:          "https://twitter.com/user4",
			Instagram:        "https://instagram.com/user4",
			JobTitle:         "Product Manager",
			Company:          "Innovative Tech",
			Website:          "https://user4.com",
		},
		{
			UserID:           users[4].ID,
			Name:             "User 5",
			ProfilePictureURL: "https://example.com/user5.jpg",
			Phone:            "9876543456",
			Gender:           models.Female,
			TheSentence:      "User 5's life motto.",
			LinkedIn:         "https://linkedin.com/in/user5",
			GitHub:           "https://github.com/user5",
			Twitter:          "https://twitter.com/user5",
			Instagram:        "https://instagram.com/user5",
			JobTitle:         "Product Manager",
			Company:          "Innovative Tech",
			Website:          "https://user5.com",
		},
	}

	for _, profile := range profiles {
		if err := tx.Create(&profile).Error; err != nil {
			rollbackAndLogError(tx, "profile", err)
			return
		}
		log.Printf("Inserted Profile: %v", profile)
	}

	tenures := []models.Tenure{
		{IsCurrentTenure: false, CurrentYear: 2025, StartDate: time.Now().AddDate(0, -1, 0), EndDate: time.Now().AddDate(0, 0, 0)},
		{IsCurrentTenure: false, CurrentYear: 2025, StartDate: time.Now().AddDate(0, -2, 0), EndDate: time.Now().AddDate(0, -1, 0)},
		{IsCurrentTenure: false, CurrentYear: 2025, StartDate: time.Now().AddDate(0, -3, 0), EndDate: time.Now().AddDate(0, -2, 0)},
		{IsCurrentTenure: false, CurrentYear: 2025, StartDate: time.Now().AddDate(0, -4, 0), EndDate: time.Now().AddDate(0, -3, 0)},
		{IsCurrentTenure: true, CurrentYear: 2026, StartDate: time.Now().AddDate(0, -4, 0), EndDate: time.Now().AddDate(0, -3, 0)},
	}

	for i := range tenures {
		if err := tx.Create(&tenures[i]).Error; err != nil {
			rollbackAndLogError(tx, "tenure", err)
			return
		}
		log.Printf("Inserted Tenure: %v", tenures[i])
	}

	memberships := []models.Membership{
		{UserID: users[0].ID, TenureID: tenures[0].ID, Role: models.Community, Position: ""},
		{UserID: users[1].ID, TenureID: tenures[1].ID, Role: models.Community, Position: ""},
		{UserID: users[2].ID, TenureID: tenures[2].ID, Role: models.Community, Position: ""},
		{UserID: users[3].ID, TenureID: tenures[4].ID, Role: models.Community, Position: ""},
	}

	for _, membership := range memberships {
		if err := tx.Create(&membership).Error; err != nil {
			rollbackAndLogError(tx, "membership", err)
			return
		}
		log.Printf("Inserted Membership: %v", membership)
	}
	
	if err := tx.Commit().Error; err != nil {
		log.Println("Error committing transaction:", err)
		tx.Rollback()
		return
	}

	log.Println("Database seeded successfully with 5 users, 5 profiles, 5 tenures, and 4 memberships!")
}

func rollbackAndLogError(tx *gorm.DB, entity string, err error) {
	tx.Rollback()
	log.Printf("Error seeding %s: %v", entity, err)
}