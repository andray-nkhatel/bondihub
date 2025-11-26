package main

import (
	"bondihub/config"
	"bondihub/services"
	"fmt"
	"os"
)

func main() {
	// Initialize config
	config.InitConfig()

	// Display config values (partially masked for security)
	cfg := config.AppConfig
	fmt.Println("Cloudinary Configuration:")
	fmt.Printf("  Cloud Name: %s\n", cfg.CloudinaryCloud)
	fmt.Printf("  API Key: %s\n", cfg.CloudinaryKey)

	// Show only first 10 chars of secret for security
	secretLen := len(cfg.CloudinarySecret)
	if secretLen > 10 {
		fmt.Printf("  API Secret: %s... (length: %d)\n", cfg.CloudinarySecret[:10], secretLen)
	} else if secretLen > 0 {
		fmt.Printf("  API Secret: %s (length: %d)\n", "***", secretLen)
	} else {
		fmt.Println("  API Secret: (EMPTY - this is the problem!)")
	}
	fmt.Println()

	// Try to initialize Cloudinary service
	fmt.Println("Attempting to initialize Cloudinary service...")
	_, err := services.NewCloudinaryService()
	if err != nil {
		fmt.Printf("❌ Failed to initialize Cloudinary: %v\n", err)
		fmt.Println("\nTo fix this:")
		fmt.Println("1. Create a .env file with your Cloudinary credentials, or")
		fmt.Println("2. Set environment variables:")
		fmt.Println("   export CLOUDINARY_CLOUD_NAME=your-cloud-name")
		fmt.Println("   export CLOUDINARY_API_KEY=your-api-key")
		fmt.Println("   export CLOUDINARY_API_SECRET=your-api-secret")
		os.Exit(1)
	}

	fmt.Println("✅ Cloudinary service initialized successfully!")
	fmt.Println()
	fmt.Println("✅ Cloudinary configuration is valid!")
	fmt.Println("   The service is ready to accept image uploads.")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("1. Restart your server to ensure it picks up the Cloudinary config")
	fmt.Println("2. Run ./test_image_upload.sh to test the full upload flow")
}
