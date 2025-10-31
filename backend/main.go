// @title BondiHub API
// @version 1.0
// @description House Renting Service API for Zambia - Find and rent houses across Zambia with secure payments via MTN MoMo and Airtel Money
// @termsOfService http://swagger.io/terms/

// @contact.name BondiHub Support
// @contact.email support@bondihub.com
// @contact.url https://bondihub.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"bondihub/config"
	"bondihub/docs"
	"bondihub/middleware"
	"bondihub/routes"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Initialize Swagger docs
	docs.SwaggerInfo.Title = "BondiHub API"
	docs.SwaggerInfo.Description = "House Renting Service API for Zambia"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Load configuration
	config.InitConfig()

	// Initialize database
	config.InitDB()
	config.AutoMigrate()

	// Set Gin mode
	gin.SetMode(config.AppConfig.GinMode)

	// Create Gin router
	r := gin.Default()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	// Setup routes
	routes.SetupRoutes(r)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	port := config.AppConfig.Port
	log.Printf("Starting BondiHub API server on port %s", port)
	log.Printf("Environment: %s", config.AppConfig.GinMode)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
