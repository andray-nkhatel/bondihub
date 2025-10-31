package handlers

import (
	"bondihub/config"
	"bondihub/models"
	"bondihub/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReviewHandler handles review-related requests
type ReviewHandler struct{}

// NewReviewHandler creates a new review handler
func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{}
}

// CreateReviewRequest represents the request structure for creating a review
type CreateReviewRequest struct {
	HouseID uuid.UUID `json:"house_id" binding:"required"`
	Rating  int       `json:"rating" binding:"required,min=1,max=5"`
	Comment string    `json:"comment" binding:"required,min=10"`
}

// CreateReview handles creating a new review
func (rh *ReviewHandler) CreateReview(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	if userModel.Role != models.RoleTenant {
		utils.ForbiddenResponse(c, "Only tenants can create reviews")
		return
	}

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Get house
	var house models.House
	if err := config.DB.First(&house, req.HouseID).Error; err != nil {
		utils.NotFoundResponse(c, "House not found")
		return
	}

	// Check if user has an active rental agreement for this house
	var agreement models.RentalAgreement
	if err := config.DB.Where("house_id = ? AND tenant_id = ? AND status = ?",
		req.HouseID, userModel.ID, models.AgreementStatusActive).First(&agreement).Error; err != nil {
		utils.ForbiddenResponse(c, "You can only review houses you have rented")
		return
	}

	// Check if user has already reviewed this house
	var existingReview models.Review
	if err := config.DB.Where("house_id = ? AND tenant_id = ?", req.HouseID, userModel.ID).First(&existingReview).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "You have already reviewed this house", nil)
		return
	}

	// Create review
	review := models.Review{
		TenantID: userModel.ID,
		HouseID:  req.HouseID,
		Rating:   req.Rating,
		Comment:  req.Comment,
	}

	if err := config.DB.Create(&review).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create review", err)
		return
	}

	// Load relationships
	config.DB.Preload("Tenant").Preload("House").First(&review, review.ID)

	// Create notification for landlord
	notification := models.Notification{
		UserID:  house.LandlordID,
		Title:   "New Review Received",
		Message: fmt.Sprintf("You received a %d-star review for %s", req.Rating, house.Title),
		Type:    "review",
	}
	config.DB.Create(&notification)

	utils.SuccessResponse(c, http.StatusCreated, "Review created successfully", gin.H{
		"review": review,
	})
}

// GetReviews handles getting reviews for a house
func (rh *ReviewHandler) GetReviews(c *gin.Context) {
	houseID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(houseID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid house ID", err)
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := config.DB.Model(&models.Review{}).
		Where("house_id = ?", id).
		Preload("Tenant")

	// Get total count
	var total int64
	query.Count(&total)

	// Get reviews
	var reviews []models.Review
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&reviews).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch reviews", err)
		return
	}

	// Calculate average rating
	var avgRating float64
	config.DB.Model(&models.Review{}).Where("house_id = ?", id).Select("AVG(rating)").Scan(&avgRating)

	// Calculate rating distribution
	var ratingDistribution []struct {
		Rating int64 `json:"rating"`
		Count  int64 `json:"count"`
	}
	config.DB.Model(&models.Review{}).
		Where("house_id = ?", id).
		Select("rating, COUNT(*) as count").
		Group("rating").
		Order("rating").
		Find(&ratingDistribution)

	// Calculate pagination info
	totalPages := (total + int64(limit) - 1) / int64(limit)

	utils.SuccessResponse(c, http.StatusOK, "Reviews retrieved successfully", gin.H{
		"reviews":             reviews,
		"average_rating":      avgRating,
		"rating_distribution": ratingDistribution,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetUserReviews handles getting reviews by a specific user
func (rh *ReviewHandler) GetUserReviews(c *gin.Context) {
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
	query := config.DB.Model(&models.Review{}).
		Where("tenant_id = ?", userModel.ID).
		Preload("House")

	// Get total count
	var total int64
	query.Count(&total)

	// Get reviews
	var reviews []models.Review
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&reviews).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch reviews", err)
		return
	}

	// Calculate pagination info
	totalPages := (total + int64(limit) - 1) / int64(limit)

	utils.SuccessResponse(c, http.StatusOK, "Reviews retrieved successfully", gin.H{
		"reviews": reviews,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// UpdateReview handles updating a review
func (rh *ReviewHandler) UpdateReview(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	reviewID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(reviewID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid review ID", err)
		return
	}

	// Get review
	var review models.Review
	if err := config.DB.First(&review, id).Error; err != nil {
		utils.NotFoundResponse(c, "Review not found")
		return
	}

	// Check if user owns the review or is admin
	if review.TenantID != userModel.ID && userModel.Role != models.RoleAdmin {
		utils.ForbiddenResponse(c, "You can only update your own reviews")
		return
	}

	var req struct {
		Rating  int    `json:"rating" binding:"min=1,max=5"`
		Comment string `json:"comment" binding:"min=10"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Update review
	if req.Rating > 0 {
		review.Rating = req.Rating
	}
	if req.Comment != "" {
		review.Comment = req.Comment
	}

	if err := config.DB.Save(&review).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update review", err)
		return
	}

	// Load relationships
	config.DB.Preload("Tenant").Preload("House").First(&review, review.ID)

	utils.SuccessResponse(c, http.StatusOK, "Review updated successfully", gin.H{
		"review": review,
	})
}

// DeleteReview handles deleting a review
func (rh *ReviewHandler) DeleteReview(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	reviewID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(reviewID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid review ID", err)
		return
	}

	// Get review
	var review models.Review
	if err := config.DB.First(&review, id).Error; err != nil {
		utils.NotFoundResponse(c, "Review not found")
		return
	}

	// Check if user owns the review or is admin
	if review.TenantID != userModel.ID && userModel.Role != models.RoleAdmin {
		utils.ForbiddenResponse(c, "You can only delete your own reviews")
		return
	}

	// Soft delete review
	if err := config.DB.Delete(&review).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete review", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Review deleted successfully", nil)
}
