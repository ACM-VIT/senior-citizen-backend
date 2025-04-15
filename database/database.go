package database

import (
	"fmt"
	"log"
	"time"

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
	} 
	var now time.Time
	DB.Raw("SELECT NOW()").Scan(&now)

	fmt.Println(now)
}
