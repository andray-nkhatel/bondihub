package main

import (
	"bondihub/config"
	"bondihub/models"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load configuration
	config.InitConfig()

	// Initialize database
	config.InitDB()
	config.AutoMigrate()

	log.Println("Starting database seeding...")

	// Create a landlord user
	landlordID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	landlord := models.User{
		ID:           landlordID,
		FullName:     "John Landlord",
		Email:        "landlord@bondihub.com",
		PasswordHash: string(hashedPassword),
		Phone:        "+260971234567",
		Role:         models.RoleLandlord,
		IsActive:     true,
		IsVerified:   true,
	}

	if err := config.DB.FirstOrCreate(&landlord, models.User{Email: landlord.Email}).Error; err != nil {
		log.Printf("Error creating landlord: %v", err)
	} else {
		log.Printf("Landlord created/found: %s", landlord.Email)
		landlordID = landlord.ID
	}

	// Create sample houses
	houses := []models.House{
		{
			LandlordID:  landlordID,
			Title:       "Modern 2 Bedroom Apartment in Kabulonga",
			Description: "Beautiful modern apartment with spacious rooms, fitted kitchen, and secure parking. Located in the heart of Kabulonga with easy access to shops and restaurants.",
			Address:     "123 Kabulonga Road, Lusaka",
			MonthlyRent: 8500,
			Status:      models.StatusAvailable,
			HouseType:   models.TypeApartment,
			Bedrooms:    2,
			Bathrooms:   2,
			Area:        85,
			Latitude:    -15.4067,
			Longitude:   28.3228,
		},
		{
			LandlordID:  landlordID,
			Title:       "Spacious 3 Bedroom House in Roma",
			Description: "Large family home with garden, garage, and servant quarters. Quiet neighborhood perfect for families. Recently renovated with modern finishes.",
			Address:     "45 Roma Park, Lusaka",
			MonthlyRent: 15000,
			Status:      models.StatusAvailable,
			HouseType:   models.TypeHouse,
			Bedrooms:    3,
			Bathrooms:   2,
			Area:        180,
			Latitude:    -15.3982,
			Longitude:   28.2891,
		},
		{
			LandlordID:  landlordID,
			Title:       "Cozy Studio in Mass Media",
			Description: "Perfect starter apartment for young professionals. Includes kitchenette, bathroom, and balcony. Walking distance to Mass Media complex.",
			Address:     "78 Mass Media, Lusaka",
			MonthlyRent: 3500,
			Status:      models.StatusAvailable,
			HouseType:   models.TypeStudio,
			Bedrooms:    1,
			Bathrooms:   1,
			Area:        35,
			Latitude:    -15.4234,
			Longitude:   28.3012,
		},
		{
			LandlordID:  landlordID,
			Title:       "Luxury 4 Bedroom Villa in Ibex Hill",
			Description: "Executive home with swimming pool, landscaped garden, and 24-hour security. Premium finishes throughout. Perfect for diplomats and executives.",
			Address:     "12 Ibex Hill, Lusaka",
			MonthlyRent: 25000,
			Status:      models.StatusAvailable,
			HouseType:   models.TypeHouse,
			Bedrooms:    4,
			Bathrooms:   3,
			Area:        350,
			Latitude:    -15.3756,
			Longitude:   28.3567,
		},
		{
			LandlordID:  landlordID,
			Title:       "2 Bedroom Townhouse in Chalala",
			Description: "Modern townhouse in gated community. Features open-plan living, fitted kitchen, and private garden. Community amenities include playground and gym.",
			Address:     "Plot 34, Chalala, Lusaka",
			MonthlyRent: 7000,
			Status:      models.StatusAvailable,
			HouseType:   models.TypeTownhouse,
			Bedrooms:    2,
			Bathrooms:   1,
			Area:        95,
			Latitude:    -15.4456,
			Longitude:   28.2678,
		},
		{
			LandlordID:  landlordID,
			Title:       "1 Bedroom Flat in Woodlands",
			Description: "Affordable flat in popular Woodlands area. Close to schools, shops, and public transport. Ideal for singles or couples.",
			Address:     "56 Woodlands, Lusaka",
			MonthlyRent: 4500,
			Status:      models.StatusAvailable,
			HouseType:   models.TypeApartment,
			Bedrooms:    1,
			Bathrooms:   1,
			Area:        50,
			Latitude:    -15.4123,
			Longitude:   28.3345,
		},
		{
			LandlordID:  landlordID,
			Title:       "3 Bedroom House in Olympia",
			Description: "Well-maintained family home with large yard. Features include veranda, domestic quarters, and carport. Quiet residential area.",
			Address:     "89 Olympia Park, Lusaka",
			MonthlyRent: 9500,
			Status:      models.StatusOccupied,
			HouseType:   models.TypeHouse,
			Bedrooms:    3,
			Bathrooms:   2,
			Area:        150,
			Latitude:    -15.3891,
			Longitude:   28.3123,
		},
		{
			LandlordID:  landlordID,
			Title:       "Executive Apartment in Longacres",
			Description: "High-end apartment with city views. Features include air conditioning, backup power, and secure parking. Building has gym and rooftop terrace.",
			Address:     "Longacres Business Park, Lusaka",
			MonthlyRent: 12000,
			Status:      models.StatusAvailable,
			HouseType:   models.TypeApartment,
			Bedrooms:    2,
			Bathrooms:   2,
			Area:        110,
			Latitude:    -15.4012,
			Longitude:   28.2956,
		},
	}

	for _, house := range houses {
		if err := config.DB.FirstOrCreate(&house, models.House{Title: house.Title}).Error; err != nil {
			log.Printf("Error creating house '%s': %v", house.Title, err)
		} else {
			log.Printf("House created/found: %s", house.Title)
		}
	}

	log.Println("Database seeding completed!")
}

