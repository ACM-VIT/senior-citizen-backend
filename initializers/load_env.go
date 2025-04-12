package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DATABASE_URL   string
	SHOULD_MIGRATE string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := Config{
		DATABASE_URL:   os.Getenv("DATABASE_URL"),
		SHOULD_MIGRATE: os.Getenv("SHOULD_MIGRATE"),
	}

	if config.DATABASE_URL == "" || config.SHOULD_MIGRATE == "" {
		log.Fatal("Missing required environment variables")
	}

	return config
}
