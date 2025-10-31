package handlers

import (
	"bondihub/config"
	"bondihub/models"
	"bondihub/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// FavoriteHandler handles favorite-related requests
type FavoriteHandler struct{}

// NewFavoriteHandler creates a new favorite handler
func NewFavoriteHandler() *FavoriteHandler {
	return &FavoriteHandler{}
}

// AddToFavorites handles adding a house to favorites
func (fh *FavoriteHandler) AddToFavorites(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	if userModel.Role != models.RoleTenant {
		utils.ForbiddenResponse(c, "Only tenants can add houses to favorites")
		return
	}

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

	// Check if already in favorites
	var existingFavorite models.Favorite
	if err := config.DB.Where("tenant_id = ? AND house_id = ?", userModel.ID, id).First(&existingFavorite).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "House already in favorites", nil)
		return
	}

	// Add to favorites
	favorite := models.Favorite{
		TenantID: userModel.ID,
		HouseID:  id,
	}

	if err := config.DB.Create(&favorite).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to add to favorites", err)
		return
	}

	// Load relationships
	config.DB.Preload("House").Preload("Tenant").First(&favorite, favorite.ID)

	utils.SuccessResponse(c, http.StatusCreated, "House added to favorites successfully", gin.H{
		"favorite": favorite,
	})
}

// RemoveFromFavorites handles removing a house from favorites
func (fh *FavoriteHandler) RemoveFromFavorites(c *gin.Context) {
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

	// Get favorite
	var favorite models.Favorite
	if err := config.DB.Where("tenant_id = ? AND house_id = ?", userModel.ID, id).First(&favorite).Error; err != nil {
		utils.NotFoundResponse(c, "House not in favorites")
		return
	}

	// Remove from favorites
	if err := config.DB.Delete(&favorite).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to remove from favorites", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "House removed from favorites successfully", nil)
}

// GetFavorites handles getting user's favorite houses
func (fh *FavoriteHandler) GetFavorites(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := config.DB.Model(&models.Favorite{}).
		Where("tenant_id = ?", userModel.ID).
		Preload("House.Landlord").
		Preload("House.Images")

	// Get total count
	var total int64
	query.Count(&total)

	// Get favorites
	var favorites []models.Favorite
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&favorites).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch favorites", err)
		return
	}

	// Calculate pagination info
	totalPages := (total + int64(limit) - 1) / int64(limit)

	utils.SuccessResponse(c, http.StatusOK, "Favorites retrieved successfully", gin.H{
		"favorites": favorites,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// CheckFavorite handles checking if a house is in favorites
func (fh *FavoriteHandler) CheckFavorite(c *gin.Context) {
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

	// Check if in favorites
	var favorite models.Favorite
	isFavorite := config.DB.Where("tenant_id = ? AND house_id = ?", userModel.ID, id).First(&favorite).Error == nil

	utils.SuccessResponse(c, http.StatusOK, "Favorite status retrieved successfully", gin.H{
		"is_favorite": isFavorite,
	})
}
