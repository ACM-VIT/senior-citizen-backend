package database

import (
	"fmt"
	"log"
	"time"

	"acm-portal/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB(cfg config.Config) {
	dsn := cfg.DatabaseURL

	var err error
	for attempt := 1; attempt <= 5; attempt++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Database connection failed (attempt %d/5): %v", attempt, err)
		time.Sleep(time.Duration(attempt) * time.Second)
	}
	if err != nil {
		log.Fatal("Failed to connect to database after retries: ", err)
	}
	var now time.Time
	DB.Raw("SELECT NOW()").Scan(&now)

	fmt.Println(now)
}
