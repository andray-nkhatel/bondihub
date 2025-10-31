package routes

import (
	"bondihub/handlers"
	"bondihub/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(r *gin.Engine) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler()
	houseHandler := handlers.NewHouseHandler()
	paymentHandler := handlers.NewPaymentHandler()
	rentalHandler := handlers.NewRentalHandler()
	reviewHandler := handlers.NewReviewHandler()
	maintenanceHandler := handlers.NewMaintenanceHandler()
	favoriteHandler := handlers.NewFavoriteHandler()
	notificationHandler := handlers.NewNotificationHandler()
	adminHandler := handlers.NewAdminHandler()

	// API version 1
	v1 := r.Group("/api/v1")

	// Public routes
	public := v1.Group("/")
	{
		// Authentication routes
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/logout", authHandler.Logout)

		// Public house routes (browse houses)
		public.GET("/houses", houseHandler.GetHouses)
		public.GET("/houses/:id", houseHandler.GetHouse)
		public.GET("/houses/:id/reviews", reviewHandler.GetReviews)
	}

	// Protected routes (require authentication)
	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// User profile routes
		protected.GET("/auth/profile", authHandler.GetProfile)
		protected.PUT("/auth/profile", authHandler.UpdateProfile)
		protected.PUT("/auth/change-password", authHandler.ChangePassword)

		// House management routes (landlords and admins)
		houses := protected.Group("/houses")
		houses.Use(middleware.LandlordOrAdminMiddleware())
		{
			houses.POST("", houseHandler.CreateHouse)
			houses.PUT("/:id", houseHandler.UpdateHouse)
			houses.DELETE("/:id", houseHandler.DeleteHouse)
			houses.POST("/:id/images", houseHandler.UploadHouseImage)
			houses.DELETE("/images/:imageId", houseHandler.DeleteHouseImage)
		}

		// Payment routes
		payments := protected.Group("/payments")
		{
			payments.POST("", paymentHandler.ProcessPayment)
			payments.GET("", paymentHandler.GetPayments)
			payments.GET("/:id", paymentHandler.GetPayment)
			payments.GET("/stats", paymentHandler.GetPaymentStats)
		}

		// Rental agreement routes
		rentals := protected.Group("/rentals")
		{
			rentals.POST("", rentalHandler.CreateRentalAgreement)
			rentals.GET("", rentalHandler.GetRentalAgreements)
			rentals.GET("/:id", rentalHandler.GetRentalAgreement)
			rentals.PUT("/:id", rentalHandler.UpdateRentalAgreement)
			rentals.PUT("/:id/terminate", rentalHandler.TerminateRentalAgreement)
		}

		// Review routes
		reviews := protected.Group("/reviews")
		{
			reviews.POST("", reviewHandler.CreateReview)
			reviews.GET("/my", reviewHandler.GetUserReviews)
			reviews.PUT("/:id", reviewHandler.UpdateReview)
			reviews.DELETE("/:id", reviewHandler.DeleteReview)
		}

		// Maintenance request routes
		maintenance := protected.Group("/maintenance")
		{
			maintenance.POST("", maintenanceHandler.CreateMaintenanceRequest)
			maintenance.GET("", maintenanceHandler.GetMaintenanceRequests)
			maintenance.GET("/:id", maintenanceHandler.GetMaintenanceRequest)
			maintenance.PUT("/:id", maintenanceHandler.UpdateMaintenanceRequest)
			maintenance.GET("/stats", maintenanceHandler.GetMaintenanceStats)
		}

		// Favorite routes (tenants only)
		favorites := protected.Group("/favorites")
		favorites.Use(middleware.TenantOrAdminMiddleware())
		{
			favorites.POST("/:id", favoriteHandler.AddToFavorites)
			favorites.DELETE("/:id", favoriteHandler.RemoveFromFavorites)
			favorites.GET("", favoriteHandler.GetFavorites)
			favorites.GET("/:id/check", favoriteHandler.CheckFavorite)
		}

		// Notification routes
		notifications := protected.Group("/notifications")
		{
			notifications.GET("", notificationHandler.GetNotifications)
			notifications.GET("/:id", notificationHandler.GetNotification)
			notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
			notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
			notifications.DELETE("/:id", notificationHandler.DeleteNotification)
			notifications.GET("/stats", notificationHandler.GetNotificationStats)
		}
	}

	// Admin routes (admin only)
	admin := v1.Group("/admin")
	admin.Use(middleware.AdminOnlyMiddleware())
	{
		admin.GET("/dashboard", adminHandler.GetDashboardStats)
		admin.GET("/users", adminHandler.GetUsers)
		admin.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
		admin.GET("/reports", adminHandler.GetReports)
	}

	// Health check route
	// @Summary Health check
	// @Description Check if the API is running
	// @Tags System
	// @Produce json
	// @Success 200 {object} map[string]interface{} "API is running"
	// @Router /health [get]
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "BondiHub API is running",
		})
	})
}
