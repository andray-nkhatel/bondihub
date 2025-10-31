package handlers

import (
	"bondihub/config"
	"bondihub/models"
	"bondihub/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RentalHandler handles rental agreement-related requests
type RentalHandler struct{}

// NewRentalHandler creates a new rental handler
func NewRentalHandler() *RentalHandler {
	return &RentalHandler{}
}

// CreateRentalAgreementRequest represents the request structure for creating a rental agreement
type CreateRentalAgreementRequest struct {
	HouseID    uuid.UUID `json:"house_id" binding:"required"`
	TenantID   uuid.UUID `json:"tenant_id" binding:"required"`
	StartDate  string    `json:"start_date" binding:"required"`
	EndDate    string    `json:"end_date" binding:"required"`
	RentAmount float64   `json:"rent_amount" binding:"required,min=0"`
	Deposit    float64   `json:"deposit" binding:"required,min=0"`
}

// CreateRentalAgreement handles creating a new rental agreement
// CreateRentalAgreement creates a new rental agreement
// @Summary Create rental agreement
// @Description Create a new rental agreement between landlord and tenant
// @Tags Rentals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateRentalAgreementRequest true "Rental agreement details"
// @Success 201 {object} map[string]interface{} "Rental agreement created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /rentals [post]
func (rh *RentalHandler) CreateRentalAgreement(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	if userModel.Role != models.RoleLandlord && userModel.Role != models.RoleAdmin {
		utils.ForbiddenResponse(c, "Only landlords can create rental agreements")
		return
	}

	var req CreateRentalAgreementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid start date format", err)
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid end date format", err)
		return
	}

	// Validate dates
	if endDate.Before(startDate) {
		utils.ErrorResponse(c, http.StatusBadRequest, "End date must be after start date", nil)
		return
	}

	// Get house
	var house models.House
	if err := config.DB.First(&house, req.HouseID).Error; err != nil {
		utils.NotFoundResponse(c, "House not found")
		return
	}

	// Check if user owns the house or is admin
	if house.LandlordID != userModel.ID && userModel.Role != models.RoleAdmin {
		utils.ForbiddenResponse(c, "You can only create agreements for your own houses")
		return
	}

	// Check if house is available
	if house.Status != models.StatusAvailable {
		utils.ErrorResponse(c, http.StatusBadRequest, "House is not available for rent", nil)
		return
	}

	// Get tenant
	var tenant models.User
	if err := config.DB.Where("id = ? AND role = ?", req.TenantID, models.RoleTenant).First(&tenant).Error; err != nil {
		utils.NotFoundResponse(c, "Tenant not found")
		return
	}

	// Check if there's already an active agreement for this house
	var existingAgreement models.RentalAgreement
	if err := config.DB.Where("house_id = ? AND status = ?", req.HouseID, models.AgreementStatusActive).First(&existingAgreement).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "House already has an active rental agreement", nil)
		return
	}

	// Create rental agreement
	agreement := models.RentalAgreement{
		HouseID:    req.HouseID,
		TenantID:   req.TenantID,
		StartDate:  startDate,
		EndDate:    endDate,
		RentAmount: req.RentAmount,
		Deposit:    req.Deposit,
		Status:     models.AgreementStatusActive,
	}

	if err := config.DB.Create(&agreement).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create rental agreement", err)
		return
	}

	// Update house status to occupied
	house.Status = models.StatusOccupied
	config.DB.Save(&house)

	// Load relationships
	config.DB.Preload("House").Preload("Tenant").First(&agreement, agreement.ID)

	// Create notifications
	tenantNotification := models.Notification{
		UserID:  req.TenantID,
		Title:   "New Rental Agreement",
		Message: fmt.Sprintf("You have a new rental agreement for %s", house.Title),
		Type:    "agreement",
	}
	config.DB.Create(&tenantNotification)

	utils.SuccessResponse(c, http.StatusCreated, "Rental agreement created successfully", gin.H{
		"agreement": agreement,
	})
}

// GetRentalAgreements handles getting rental agreements for a user
func (rh *RentalHandler) GetRentalAgreements(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := config.DB.Model(&models.RentalAgreement{}).Preload("House").Preload("Tenant")

	// Apply filters based on user role
	if userModel.Role == models.RoleTenant {
		query = query.Where("tenant_id = ?", userModel.ID)
	} else if userModel.Role == models.RoleLandlord {
		query = query.Joins("JOIN houses ON rental_agreements.house_id = houses.id").
			Where("houses.landlord_id = ?", userModel.ID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get agreements
	var agreements []models.RentalAgreement
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&agreements).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch rental agreements", err)
		return
	}

	// Calculate pagination info
	totalPages := (total + int64(limit) - 1) / int64(limit)

	utils.SuccessResponse(c, http.StatusOK, "Rental agreements retrieved successfully", gin.H{
		"agreements": agreements,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetRentalAgreement handles getting a single rental agreement by ID
func (rh *RentalHandler) GetRentalAgreement(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	agreementID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(agreementID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid agreement ID", err)
		return
	}

	var agreement models.RentalAgreement
	if err := config.DB.Preload("House").Preload("Tenant").Preload("Payments").First(&agreement, id).Error; err != nil {
		utils.NotFoundResponse(c, "Rental agreement not found")
		return
	}

	// Check if user has access to this agreement
	hasAccess := false
	if userModel.Role == models.RoleAdmin {
		hasAccess = true
	} else if userModel.Role == models.RoleTenant && agreement.TenantID == userModel.ID {
		hasAccess = true
	} else if userModel.Role == models.RoleLandlord && agreement.House.LandlordID == userModel.ID {
		hasAccess = true
	}

	if !hasAccess {
		utils.ForbiddenResponse(c, "You don't have access to this agreement")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Rental agreement retrieved successfully", gin.H{
		"agreement": agreement,
	})
}

// UpdateRentalAgreement handles updating a rental agreement
func (rh *RentalHandler) UpdateRentalAgreement(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	agreementID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(agreementID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid agreement ID", err)
		return
	}

	// Get agreement
	var agreement models.RentalAgreement
	if err := config.DB.Preload("House").First(&agreement, id).Error; err != nil {
		utils.NotFoundResponse(c, "Rental agreement not found")
		return
	}

	// Check if user has access to this agreement
	hasAccess := false
	if userModel.Role == models.RoleAdmin {
		hasAccess = true
	} else if agreement.House.LandlordID == userModel.ID {
		hasAccess = true
	}

	if !hasAccess {
		utils.ForbiddenResponse(c, "You don't have access to this agreement")
		return
	}

	var req struct {
		Status string `json:"status" binding:"oneof=active terminated expired"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Update agreement status
	agreement.Status = models.AgreementStatus(req.Status)
	agreement.UpdatedAt = time.Now()

	if err := config.DB.Save(&agreement).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update rental agreement", err)
		return
	}

	// If agreement is terminated or expired, update house status to available
	if agreement.Status == models.AgreementStatusTerminated || agreement.Status == models.AgreementStatusExpired {
		var house models.House
		config.DB.First(&house, agreement.HouseID)
		house.Status = models.StatusAvailable
		config.DB.Save(&house)
	}

	// Load relationships
	config.DB.Preload("House").Preload("Tenant").First(&agreement, agreement.ID)

	utils.SuccessResponse(c, http.StatusOK, "Rental agreement updated successfully", gin.H{
		"agreement": agreement,
	})
}

// TerminateRentalAgreement handles terminating a rental agreement
func (rh *RentalHandler) TerminateRentalAgreement(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	agreementID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(agreementID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid agreement ID", err)
		return
	}

	// Get agreement
	var agreement models.RentalAgreement
	if err := config.DB.Preload("House").First(&agreement, id).Error; err != nil {
		utils.NotFoundResponse(c, "Rental agreement not found")
		return
	}

	// Check if user has access to this agreement
	hasAccess := false
	if userModel.Role == models.RoleAdmin {
		hasAccess = true
	} else if agreement.House.LandlordID == userModel.ID {
		hasAccess = true
	}

	if !hasAccess {
		utils.ForbiddenResponse(c, "You don't have access to this agreement")
		return
	}

	// Check if agreement is active
	if agreement.Status != models.AgreementStatusActive {
		utils.ErrorResponse(c, http.StatusBadRequest, "Only active agreements can be terminated", nil)
		return
	}

	// Terminate agreement
	agreement.Status = models.AgreementStatusTerminated
	agreement.UpdatedAt = time.Now()

	if err := config.DB.Save(&agreement).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to terminate rental agreement", err)
		return
	}

	// Update house status to available
	var house models.House
	config.DB.First(&house, agreement.HouseID)
	house.Status = models.StatusAvailable
	config.DB.Save(&house)

	// Create notification for tenant
	notification := models.Notification{
		UserID:  agreement.TenantID,
		Title:   "Rental Agreement Terminated",
		Message: fmt.Sprintf("Your rental agreement for %s has been terminated", house.Title),
		Type:    "agreement",
	}
	config.DB.Create(&notification)

	utils.SuccessResponse(c, http.StatusOK, "Rental agreement terminated successfully", gin.H{
		"agreement": agreement,
	})
}
