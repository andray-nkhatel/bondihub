package config

import (
	"fmt"
	"log"

	"bondihub/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
	var err error

	// Database connection string with defaults
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "bondihub"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)

	// Connect to database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
}

// AutoMigrate runs database migrations
func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.House{},
		&models.HouseImage{},
		&models.RentalAgreement{},
		&models.Payment{},
		&models.Review{},
		&models.MaintenanceRequest{},
		&models.Favorite{},
		&models.Notification{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed successfully")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
