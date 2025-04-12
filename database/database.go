package database

import (
	"log"

	"acm-portal/initializers"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	dsn := initializers.LoadConfig().DATABASE_URL

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	} else {
		log.Println("Connected to database")
	}
}
