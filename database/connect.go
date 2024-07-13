package database

import (
	"fmt"
	"log"
	"notes_api/config"
	"notes_api/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// connection for db creds
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Config("DB_HOST"), config.Config("DB_PORT"), config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"), config.Config("DB_SSLMODE"),
	)

	// connection for db url
	// dsn := config.Config("DB_URL")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	} else {
		fmt.Println("Connected to database")
	}

	err = DB.AutoMigrate(&model.User{}, &model.Note{})
	if err != nil {
		fmt.Println("Error occured while migrating", err)
	} else {
		fmt.Println("Migrated models successfully")
	}
}