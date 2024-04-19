package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Andrewalifb/games-rental-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB


// func DatabaseConnection() {


// 	host := os.Getenv("HOST")
// 	port := os.Getenv("PORT")
// 	dbName := os.Getenv("DB_NAME")
// 	dbUser := os.Getenv("DB_USER")
// 	password := os.Getenv("PASSWORD")
// 	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
// 		host,
// 		port,
// 		dbUser,
// 		dbName,
// 		password,
// 	)

// 	DB, dbError = gorm.Open(postgres.Open(dsn), &gorm.Config{})

// 	if dbError != nil {
// 		log.Fatal("Error connecting to the database...", dbError)
// 	}
// 	log.Println("Database connection successful...")
// }

func DatabaseConnection() {
	// Use the DATABASE_URL from Heroku config vars
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println("DATABASE_URL environment variable not set")
		return
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		fmt.Println("There is an error while connecting to the database ", err)
		panic(err)
	} else {
		DB = db

		fmt.Println("Successfully connected to database!")
	}
}

func Migrate() {
	DB.AutoMigrate(&models.User{}, &models.PaymentHistory{}, &models.GameCategory{}, &models.GamePlatform{}, &models.Game{}, &models.Cart{}, &models.RentalTransaction{}, &models.RentMaintenance{})
  // Check if the sequence exists
	var count int64
	DB.Raw("SELECT COUNT(*) FROM pg_sequences WHERE sequencename = 'rental_id_seq'").Scan(&count)
	// If the sequence does not exist, create it
	if count == 0 {
	DB.Exec("CREATE SEQUENCE rental_id_seq START 1000")
	}
	
	log.Println("Database Migration Completed!")
}