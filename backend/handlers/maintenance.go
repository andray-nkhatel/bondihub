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

// MaintenanceHandler handles maintenance request-related requests
type MaintenanceHandler struct{}

// NewMaintenanceHandler creates a new maintenance handler
func NewMaintenanceHandler() *MaintenanceHandler {
	return &MaintenanceHandler{}
}

// CreateMaintenanceRequest represents the request structure for creating a maintenance request
type CreateMaintenanceRequest struct {
	HouseID     uuid.UUID `json:"house_id" binding:"required"`
	Title       string    `json:"title" binding:"required,min=5,max=200"`
	Description string    `json:"description" binding:"required,min=10"`
	Priority    string    `json:"priority" binding:"oneof=low medium high urgent"`
}

// CreateMaintenanceRequest handles creating a new maintenance request
func (mh *MaintenanceHandler) CreateMaintenanceRequest(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	if userModel.Role != models.RoleTenant {
		utils.ForbiddenResponse(c, "Only tenants can create maintenance requests")
		return
	}

	var req CreateMaintenanceRequest
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
		utils.ForbiddenResponse(c, "You can only create maintenance requests for houses you rent")
		return
	}

	// Create maintenance request
	maintenanceRequest := models.MaintenanceRequest{
		TenantID:    userModel.ID,
		HouseID:     req.HouseID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      models.MaintenanceStatusPending,
		ReportedAt:  time.Now(),
	}

	if err := config.DB.Create(&maintenanceRequest).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create maintenance request", err)
		return
	}

	// Load relationships
	config.DB.Preload("Tenant").Preload("House").First(&maintenanceRequest, maintenanceRequest.ID)

	// Create notification for landlord
	notification := models.Notification{
		UserID:  house.LandlordID,
		Title:   "New Maintenance Request",
		Message: fmt.Sprintf("New maintenance request for %s: %s", house.Title, req.Title),
		Type:    "maintenance",
	}
	config.DB.Create(&notification)

	utils.SuccessResponse(c, http.StatusCreated, "Maintenance request created successfully", gin.H{
		"maintenance_request": maintenanceRequest,
	})
}

// GetMaintenanceRequests handles getting maintenance requests
func (mh *MaintenanceHandler) GetMaintenanceRequests(c *gin.Context) {
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
	priority := c.Query("priority")

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := config.DB.Model(&models.MaintenanceRequest{}).Preload("Tenant").Preload("House")

	// Apply filters based on user role
	if userModel.Role == models.RoleTenant {
		query = query.Where("tenant_id = ?", userModel.ID)
	} else if userModel.Role == models.RoleLandlord {
		query = query.Joins("JOIN houses ON maintenance_requests.house_id = houses.id").
			Where("houses.landlord_id = ?", userModel.ID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get maintenance requests
	var requests []models.MaintenanceRequest
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&requests).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch maintenance requests", err)
		return
	}

	// Calculate pagination info
	totalPages := (total + int64(limit) - 1) / int64(limit)

	utils.SuccessResponse(c, http.StatusOK, "Maintenance requests retrieved successfully", gin.H{
		"maintenance_requests": requests,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetMaintenanceRequest handles getting a single maintenance request by ID
func (mh *MaintenanceHandler) GetMaintenanceRequest(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	requestID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(requestID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request ID", err)
		return
	}

	var request models.MaintenanceRequest
	if err := config.DB.Preload("Tenant").Preload("House").First(&request, id).Error; err != nil {
		utils.NotFoundResponse(c, "Maintenance request not found")
		return
	}

	// Check if user has access to this request
	hasAccess := false
	if userModel.Role == models.RoleAdmin {
		hasAccess = true
	} else if userModel.Role == models.RoleTenant && request.TenantID == userModel.ID {
		hasAccess = true
	} else if userModel.Role == models.RoleLandlord && request.House.LandlordID == userModel.ID {
		hasAccess = true
	}

	if !hasAccess {
		utils.ForbiddenResponse(c, "You don't have access to this maintenance request")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Maintenance request retrieved successfully", gin.H{
		"maintenance_request": request,
	})
}

// UpdateMaintenanceRequest handles updating a maintenance request
func (mh *MaintenanceHandler) UpdateMaintenanceRequest(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	requestID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(requestID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request ID", err)
		return
	}

	// Get request
	var request models.MaintenanceRequest
	if err := config.DB.Preload("House").First(&request, id).Error; err != nil {
		utils.NotFoundResponse(c, "Maintenance request not found")
		return
	}

	// Check if user has access to this request
	hasAccess := false
	if userModel.Role == models.RoleAdmin {
		hasAccess = true
	} else if request.House.LandlordID == userModel.ID {
		hasAccess = true
	}

	if !hasAccess {
		utils.ForbiddenResponse(c, "You don't have access to this maintenance request")
		return
	}

	var req struct {
		Status string `json:"status" binding:"oneof=pending in_progress resolved cancelled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Update request status
	request.Status = models.MaintenanceRequestStatus(req.Status)
	request.UpdatedAt = time.Now()

	// Set resolved_at if status is resolved
	if request.Status == models.MaintenanceStatusResolved {
		now := time.Now()
		request.ResolvedAt = &now
	}

	if err := config.DB.Save(&request).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update maintenance request", err)
		return
	}

	// Load relationships
	config.DB.Preload("Tenant").Preload("House").First(&request, request.ID)

	// Create notification for tenant
	notification := models.Notification{
		UserID:  request.TenantID,
		Title:   "Maintenance Request Updated",
		Message: fmt.Sprintf("Your maintenance request for %s has been updated to: %s", request.House.Title, req.Status),
		Type:    "maintenance",
	}
	config.DB.Create(&notification)

	utils.SuccessResponse(c, http.StatusOK, "Maintenance request updated successfully", gin.H{
		"maintenance_request": request,
	})
}

// GetMaintenanceStats handles getting maintenance statistics
func (mh *MaintenanceHandler) GetMaintenanceStats(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)

	// Build base query
	query := config.DB.Model(&models.MaintenanceRequest{})

	// Apply filters based on user role
	if userModel.Role == models.RoleTenant {
		query = query.Where("tenant_id = ?", userModel.ID)
	} else if userModel.Role == models.RoleLandlord {
		query = query.Joins("JOIN houses ON maintenance_requests.house_id = houses.id").
			Where("houses.landlord_id = ?", userModel.ID)
	}

	// Get total requests
	var totalRequests int64
	query.Count(&totalRequests)

	// Get requests by status
	var requestsByStatus []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	query.Select("status, COUNT(*) as count").
		Group("status").
		Find(&requestsByStatus)

	// Get requests by priority
	var requestsByPriority []struct {
		Priority string `json:"priority"`
		Count    int64  `json:"count"`
	}
	query.Select("priority, COUNT(*) as count").
		Group("priority").
		Find(&requestsByPriority)

	// Get average resolution time (in days)
	var avgResolutionDays float64
	query.Where("status = ? AND resolved_at IS NOT NULL", models.MaintenanceStatusResolved).
		Select("AVG(EXTRACT(EPOCH FROM (resolved_at - reported_at))/86400)").
		Scan(&avgResolutionDays)

	utils.SuccessResponse(c, http.StatusOK, "Maintenance statistics retrieved successfully", gin.H{
		"total_requests":       totalRequests,
		"requests_by_status":   requestsByStatus,
		"requests_by_priority": requestsByPriority,
		"avg_resolution_days":  avgResolutionDays,
	})
}
