package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	JWTSecret          string
	JWTExpiresIn       time.Duration
	Port               string
	GinMode            string
	CloudinaryCloud    string
	CloudinaryKey      string
	CloudinarySecret   string
	MTNMoMoAPIURL      string
	MTNMoMoAPIKey      string
	MTNMoMoSubKey      string
	AirtelAPIURL       string
	AirtelClientID     string
	AirtelClientSecret string
	CommissionRate     float64
	FeaturedPrice      float64
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Parse JWT expiration duration
	jwtExpiresIn, err := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))
	if err != nil {
		log.Fatal("Invalid JWT_EXPIRES_IN format:", err)
	}

	// Parse commission rate
	commissionRate, err := strconv.ParseFloat(getEnv("COMMISSION_RATE", "0.05"), 64)
	if err != nil {
		log.Fatal("Invalid COMMISSION_RATE format:", err)
	}

	// Parse featured listing price
	featuredPrice, err := strconv.ParseFloat(getEnv("FEATURED_LISTING_PRICE", "500.00"), 64)
	if err != nil {
		log.Fatal("Invalid FEATURED_LISTING_PRICE format:", err)
	}

	return &Config{
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "password"),
		DBName:             getEnv("DB_NAME", "bondihub"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		JWTSecret:          getEnv("JWT_SECRET", "your-super-secret-jwt-key-here"),
		JWTExpiresIn:       jwtExpiresIn,
		Port:               getEnv("PORT", "8080"),
		GinMode:            getEnv("GIN_MODE", "debug"),
		CloudinaryCloud:    getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryKey:      getEnv("CLOUDINARY_API_KEY", ""),
		CloudinarySecret:   getEnv("CLOUDINARY_API_SECRET", ""),
		MTNMoMoAPIURL:      getEnv("MTN_MOMO_API_URL", ""),
		MTNMoMoAPIKey:      getEnv("MTN_MOMO_API_KEY", ""),
		MTNMoMoSubKey:      getEnv("MTN_MOMO_SUBSCRIPTION_KEY", ""),
		AirtelAPIURL:       getEnv("AIRTEL_MONEY_API_URL", ""),
		AirtelClientID:     getEnv("AIRTEL_MONEY_CLIENT_ID", ""),
		AirtelClientSecret: getEnv("AIRTEL_MONEY_CLIENT_SECRET", ""),
		CommissionRate:     commissionRate,
		FeaturedPrice:      featuredPrice,
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// AppConfig holds the global configuration
var AppConfig *Config

// InitConfig initializes the global configuration
func InitConfig() {
	AppConfig = Load()
}
