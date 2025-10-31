package db_utils

import (
	"log"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

func ConnectAndMigrate() (*gorm.DB, error) {
	err := godotenv.Load()
    if err != nil {
        // This is useful for production where you rely only on system environment variables
        log.Println("Note: Could not find .env file, assuming environment variables are set globally.")
    }
	// 1. Database Setup
	// NOTE: Hardcoded DSN for simplicity, but ideally this would come from environment variables.
	dsn := os.Getenv("DATABASE_URL")
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Println("Successfully connected to database!")

	// 2. AutoMigrate all models
	// Since this function is in the db_utils package, we don't need a package prefix.
	err = db.AutoMigrate(
		&User{}, 
		&Profile{}, 
		&Skill{},
		&Achievement{},
		&Injury{},
		&SocialLink{},
		&ClubProfile{},
		&SeasonStat{},
	)
	if err != nil {
		return nil, err
	}
	log.Println("Database migration completed successfully!")

	return db, nil
}