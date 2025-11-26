package handlers

import (
	"bondihub/config"
	"bondihub/models"
	"bondihub/services"
	"bondihub/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// HouseHandler handles house-related requests
type HouseHandler struct {
	cloudinaryService *services.CloudinaryService
}

// NewHouseHandler creates a new house handler
func NewHouseHandler() *HouseHandler {
	// Log config values for debugging (mask secret)
	cfg := config.AppConfig
	if cfg.CloudinaryURL != "" {
		// Mask the secret in the URL
		urlDisplay := cfg.CloudinaryURL
		if len(urlDisplay) > 30 {
			urlDisplay = urlDisplay[:20] + "***" + urlDisplay[len(urlDisplay)-10:]
		}
		log.Printf("Initializing Cloudinary service from CLOUDINARY_URL: %s", urlDisplay)
	} else {
		secretDisplay := "***"
		if len(cfg.CloudinarySecret) > 10 {
			secretDisplay = cfg.CloudinarySecret[:10] + "..."
		}
		log.Printf("Initializing Cloudinary service with Cloud: %s, Key: %s, Secret: %s...", 
			cfg.CloudinaryCloud, cfg.CloudinaryKey, secretDisplay)
	}
	
	cloudinaryService, err := services.NewCloudinaryService()
	if err != nil {
		log.Printf("❌ ERROR: Failed to initialize Cloudinary service: %v", err)
		log.Printf("   Cloud: %s, Key: %s, Secret length: %d", 
			cfg.CloudinaryCloud, cfg.CloudinaryKey, len(cfg.CloudinarySecret))
		log.Println("   Image uploads will not work until Cloudinary is properly configured")
		// Continue without Cloudinary service - uploads will fail gracefully
	} else {
		log.Println("✅ Cloudinary service initialized successfully")
	}
	return &HouseHandler{
		cloudinaryService: cloudinaryService,
	}
}

// CreateHouseRequest represents the request structure for creating a house
type CreateHouseRequest struct {
	Title       string   `json:"title" binding:"required,min=5,max=200"`
	Description string   `json:"description" binding:"required,min=10"`
	Address     string   `json:"address" binding:"required,min=10"`
	MonthlyRent float64  `json:"monthly_rent" binding:"required,min=0"`
	HouseType   string   `json:"house_type" binding:"required,oneof=apartment house studio townhouse commercial"`
	Latitude    *float64 `json:"latitude" binding:"omitempty,min=-90,max=90"`
	Longitude   *float64 `json:"longitude" binding:"omitempty,min=-180,max=180"`
	Bedrooms    int      `json:"bedrooms" binding:"min=0"`
	Bathrooms   int      `json:"bathrooms" binding:"min=0"`
	Area        float64  `json:"area" binding:"min=0"`
	IsFeatured  bool     `json:"is_featured"`
}

// UpdateHouseRequest represents the request structure for updating a house
type UpdateHouseRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Address     string  `json:"address"`
	MonthlyRent float64 `json:"monthly_rent"`
	Status      string  `json:"status"`
	HouseType   string  `json:"house_type"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Bedrooms    int     `json:"bedrooms"`
	Bathrooms   int     `json:"bathrooms"`
	Area        float64 `json:"area"`
	IsFeatured  bool    `json:"is_featured"`
}

// CreateHouse handles creating a new house
// @Summary Create house
// @Description Create a new house listing for rent (landlords and admins only)
// @Tags Houses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateHouseRequest true "House creation data"
// @Success 201 {object} map[string]interface{} "House created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Only landlords can create houses"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /houses [post]
func (hh *HouseHandler) CreateHouse(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	if userModel.Role != models.RoleLandlord && userModel.Role != models.RoleAdmin {
		utils.ForbiddenResponse(c, "Only landlords can create houses")
		return
	}

	var req CreateHouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Set default coordinates if not provided (default to 0, 0)
	// This allows users to omit coordinates, which will default to (0, 0)
	latitude := 0.0
	longitude := 0.0
	if req.Latitude != nil {
		latitude = *req.Latitude
	}
	if req.Longitude != nil {
		longitude = *req.Longitude
	}

	// Only validate coordinates if BOTH are explicitly provided (not nil)
	// If both are provided and equal to (0, 0), reject as invalid
	// If omitted, defaulting to (0, 0) is allowed
	if req.Latitude != nil && req.Longitude != nil && latitude == 0 && longitude == 0 {
		utils.ValidationErrorResponse(c, "Invalid request data", map[string]string{
			"error": "Coordinates (0, 0) are not valid. Please provide valid latitude and longitude values, or omit both fields to use default values.",
		})
		return
	}

	// Create house
	house := models.House{
		LandlordID:  userModel.ID,
		Title:       req.Title,
		Description: req.Description,
		Address:     req.Address,
		MonthlyRent: req.MonthlyRent,
		Status:      models.StatusAvailable,
		HouseType:   models.HouseType(req.HouseType),
		Latitude:    latitude,
		Longitude:   longitude,
		Bedrooms:    req.Bedrooms,
		Bathrooms:   req.Bathrooms,
		Area:        req.Area,
		IsFeatured:  req.IsFeatured,
	}

	// Set featured expiry if house is featured
	if req.IsFeatured {
		featuredUntil := time.Now().Add(30 * 24 * time.Hour) // 30 days
		house.FeaturedUntil = &featuredUntil
	}

	if err := config.DB.Create(&house).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create house", err)
		return
	}

	// Load landlord information
	config.DB.Preload("Landlord").First(&house, house.ID)

	utils.SuccessResponse(c, http.StatusCreated, "House created successfully", gin.H{
		"house": house,
	})
}

// GetHouses handles getting all houses with pagination and filters
// GetHouses retrieves a list of houses with filtering and pagination
// @Summary Get houses
// @Description Get a paginated list of houses with optional filtering
// @Tags Houses
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param house_type query string false "House type filter"
// @Param status query string false "Status filter"
// @Param min_rent query number false "Minimum rent filter"
// @Param max_rent query number false "Maximum rent filter"
// @Param bedrooms query int false "Number of bedrooms filter"
// @Param bathrooms query int false "Number of bathrooms filter"
// @Param search query string false "Search term"
// @Success 200 {object} map[string]interface{} "Houses retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /houses [get]
func (hh *HouseHandler) GetHouses(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	houseType := c.Query("house_type")
	status := c.Query("status")
	minRent, _ := strconv.ParseFloat(c.Query("min_rent"), 64)
	maxRent, _ := strconv.ParseFloat(c.Query("max_rent"), 64)
	bedrooms, _ := strconv.Atoi(c.Query("bedrooms"))
	bathrooms, _ := strconv.Atoi(c.Query("bathrooms"))
	featured := c.Query("featured") == "true"
	search := c.Query("search")

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := config.DB.Model(&models.House{}).Preload("Landlord").Preload("Images")

	// Apply filters
	if houseType != "" {
		query = query.Where("house_type = ?", houseType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if minRent > 0 {
		query = query.Where("monthly_rent >= ?", minRent)
	}
	if maxRent > 0 {
		query = query.Where("monthly_rent <= ?", maxRent)
	}
	if bedrooms > 0 {
		query = query.Where("bedrooms >= ?", bedrooms)
	}
	if bathrooms > 0 {
		query = query.Where("bathrooms >= ?", bathrooms)
	}
	if featured {
		query = query.Where("is_featured = ? AND (featured_until IS NULL OR featured_until > ?)", true, time.Now())
	}
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ? OR address ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get houses
	var houses []models.House
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&houses).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch houses", err)
		return
	}

	// Calculate pagination info
	totalPages := (total + int64(limit) - 1) / int64(limit)

	utils.SuccessResponse(c, http.StatusOK, "Houses retrieved successfully", gin.H{
		"houses": houses,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetHouse handles getting a single house by ID
// GetHouse retrieves a specific house by ID
// @Summary Get house by ID
// @Description Get detailed information about a specific house
// @Tags Houses
// @Accept json
// @Produce json
// @Param id path string true "House ID"
// @Success 200 {object} map[string]interface{} "House retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid house ID"
// @Failure 404 {object} map[string]interface{} "House not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /houses/{id} [get]
func (hh *HouseHandler) GetHouse(c *gin.Context) {
	houseID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(houseID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid house ID", err)
		return
	}

	var house models.House
	if err := config.DB.Preload("Landlord").Preload("Images").Preload("Reviews.Tenant").First(&house, id).Error; err != nil {
		utils.NotFoundResponse(c, "House not found")
		return
	}

	// Calculate average rating
	var avgRating float64
	config.DB.Model(&models.Review{}).Where("house_id = ?", house.ID).Select("AVG(rating)").Scan(&avgRating)
	house.Reviews = []models.Review{} // Clear reviews to avoid circular reference

	utils.SuccessResponse(c, http.StatusOK, "House retrieved successfully", gin.H{
		"house":          house,
		"average_rating": avgRating,
	})
}

// UpdateHouse handles updating a house
// @Summary Update house
// @Description Update an existing house listing (owner or admin only)
// @Tags Houses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "House ID"
// @Param request body UpdateHouseRequest true "House update data"
// @Success 200 {object} map[string]interface{} "House updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - You can only update your own houses"
// @Failure 404 {object} map[string]interface{} "House not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /houses/{id} [put]
func (hh *HouseHandler) UpdateHouse(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	houseID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(houseID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid house ID", err)
		return
	}

	// Get house
	var house models.House
	if err := config.DB.First(&house, id).Error; err != nil {
		utils.NotFoundResponse(c, "House not found")
		return
	}

	// Check if user owns the house or is admin
	if house.LandlordID != userModel.ID && userModel.Role != models.RoleAdmin {
		utils.ForbiddenResponse(c, "You can only update your own houses")
		return
	}

	var req UpdateHouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Update house fields
	if req.Title != "" {
		house.Title = req.Title
	}
	if req.Description != "" {
		house.Description = req.Description
	}
	if req.Address != "" {
		house.Address = req.Address
	}
	if req.MonthlyRent > 0 {
		house.MonthlyRent = req.MonthlyRent
	}
	if req.Status != "" {
		house.Status = models.HouseStatus(req.Status)
	}
	if req.HouseType != "" {
		house.HouseType = models.HouseType(req.HouseType)
	}
	if req.Latitude != 0 {
		house.Latitude = req.Latitude
	}
	if req.Longitude != 0 {
		house.Longitude = req.Longitude
	}
	if req.Bedrooms >= 0 {
		house.Bedrooms = req.Bedrooms
	}
	if req.Bathrooms >= 0 {
		house.Bathrooms = req.Bathrooms
	}
	if req.Area >= 0 {
		house.Area = req.Area
	}
	house.IsFeatured = req.IsFeatured

	// Update featured expiry if house is featured
	if req.IsFeatured && house.FeaturedUntil == nil {
		featuredUntil := time.Now().Add(30 * 24 * time.Hour) // 30 days
		house.FeaturedUntil = &featuredUntil
	}

	house.UpdatedAt = time.Now()

	if err := config.DB.Save(&house).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update house", err)
		return
	}

	// Load landlord information
	config.DB.Preload("Landlord").Preload("Images").First(&house, house.ID)

	utils.SuccessResponse(c, http.StatusOK, "House updated successfully", gin.H{
		"house": house,
	})
}

// DeleteHouse handles deleting a house
// @Summary Delete house
// @Description Delete a house listing (owner or admin only)
// @Tags Houses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "House ID"
// @Success 200 {object} map[string]interface{} "House deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid house ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - You can only delete your own houses"
// @Failure 404 {object} map[string]interface{} "House not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /houses/{id} [delete]
func (hh *HouseHandler) DeleteHouse(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	houseID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(houseID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid house ID", err)
		return
	}

	// Get house
	var house models.House
	if err := config.DB.First(&house, id).Error; err != nil {
		utils.NotFoundResponse(c, "House not found")
		return
	}

	// Check if user owns the house or is admin
	if house.LandlordID != userModel.ID && userModel.Role != models.RoleAdmin {
		utils.ForbiddenResponse(c, "You can only delete your own houses")
		return
	}

	// Soft delete house
	if err := config.DB.Delete(&house).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete house", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "House deleted successfully", nil)
}

// UploadHouseImage handles uploading images for a house
// @Summary Upload house image
// @Description Upload an image for a house listing (owner or admin only)
// @Tags Houses
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "House ID"
// @Param image formData file true "House image file"
// @Success 201 {object} map[string]interface{} "Image uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid house ID or no image file provided"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - You can only upload images for your own houses"
// @Failure 404 {object} map[string]interface{} "House not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /houses/{id}/images [post]
func (hh *HouseHandler) UploadHouseImage(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	houseID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(houseID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid house ID", err)
		return
	}

	// Get house
	var house models.House
	if err := config.DB.First(&house, id).Error; err != nil {
		utils.NotFoundResponse(c, "House not found")
		return
	}

	// Check if user owns the house or is admin
	if house.LandlordID != userModel.ID && userModel.Role != models.RoleAdmin {
		utils.ForbiddenResponse(c, "You can only upload images for your own houses")
		return
	}

	// Get uploaded file
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "No image file provided", err)
		return
	}
	defer file.Close()

	// Check if Cloudinary service is available
	if hh.cloudinaryService == nil {
		utils.InternalServerErrorResponse(c, "Image upload service is not configured", nil)
		return
	}

	// Upload to Cloudinary
	result, err := hh.cloudinaryService.UploadImage(c.Request.Context(), file, "bondihub/houses")
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to upload image", err)
		return
	}

	// Create house image record
	houseImage := models.HouseImage{
		HouseID:   house.ID,
		ImageURL:  result.SecureURL,
		IsPrimary: false,
	}

	if err := config.DB.Create(&houseImage).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to save image record", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Image uploaded successfully", gin.H{
		"image": houseImage,
	})
}

// DeleteHouseImage handles deleting a house image
// @Summary Delete house image
// @Description Delete an image from a house listing (owner or admin only)
// @Tags Houses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param imageId path string true "Image ID"
// @Success 200 {object} map[string]interface{} "Image deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid image ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - You can only delete images for your own houses"
// @Failure 404 {object} map[string]interface{} "Image not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /images/{imageId} [delete]
func (hh *HouseHandler) DeleteHouseImage(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	imageID := c.Param("imageId")

	// Parse UUID
	id, err := uuid.Parse(imageID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid image ID", err)
		return
	}

	// Get image
	var image models.HouseImage
	if err := config.DB.Preload("House").First(&image, id).Error; err != nil {
		utils.NotFoundResponse(c, "Image not found")
		return
	}

	// Check if user owns the house or is admin
	if image.House.LandlordID != userModel.ID && userModel.Role != models.RoleAdmin {
		utils.ForbiddenResponse(c, "You can only delete images for your own houses")
		return
	}

	// Delete from Cloudinary
	// Note: You might want to extract public_id from the URL
	// For now, we'll just delete from database
	if err := config.DB.Delete(&image).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete image", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Image deleted successfully", nil)
}
