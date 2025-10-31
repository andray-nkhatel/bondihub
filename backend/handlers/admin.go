package handlers

import (
	"bondihub/config"
	"bondihub/models"
	"bondihub/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AdminHandler handles admin-related requests
type AdminHandler struct{}

// NewAdminHandler creates a new admin handler
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// GetDashboardStats handles getting admin dashboard statistics
// @Summary Get admin dashboard statistics
// @Description Get comprehensive statistics for the admin dashboard
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Dashboard statistics retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin access required"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/dashboard [get]
func (ah *AdminHandler) GetDashboardStats(c *gin.Context) {
	// Get total users
	var totalUsers int64
	config.DB.Model(&models.User{}).Count(&totalUsers)

	// Get users by role
	var usersByRole []struct {
		Role  string `json:"role"`
		Count int64  `json:"count"`
	}
	config.DB.Model(&models.User{}).
		Select("role, COUNT(*) as count").
		Group("role").
		Find(&usersByRole)

	// Get total houses
	var totalHouses int64
	config.DB.Model(&models.House{}).Count(&totalHouses)

	// Get houses by status
	var housesByStatus []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	config.DB.Model(&models.House{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&housesByStatus)

	// Get total rental agreements
	var totalAgreements int64
	config.DB.Model(&models.RentalAgreement{}).Count(&totalAgreements)

	// Get active rental agreements
	var activeAgreements int64
	config.DB.Model(&models.RentalAgreement{}).Where("status = ?", models.AgreementStatusActive).Count(&activeAgreements)

	// Get total payments
	var totalPayments int64
	config.DB.Model(&models.Payment{}).Count(&totalPayments)

	// Get total revenue
	var totalRevenue float64
	config.DB.Model(&models.Payment{}).Where("status = ?", models.PaymentStatusCompleted).Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)

	// Get total commission
	var totalCommission float64
	config.DB.Model(&models.Payment{}).Where("status = ?", models.PaymentStatusCompleted).Select("COALESCE(SUM(commission), 0)").Scan(&totalCommission)

	// Get payments by method
	var paymentsByMethod []struct {
		Method string  `json:"method"`
		Count  int64   `json:"count"`
		Amount float64 `json:"amount"`
	}
	config.DB.Model(&models.Payment{}).
		Select("method, COUNT(*) as count, COALESCE(SUM(amount), 0) as amount").
		Group("method").
		Find(&paymentsByMethod)

	// Get total maintenance requests
	var totalMaintenanceRequests int64
	config.DB.Model(&models.MaintenanceRequest{}).Count(&totalMaintenanceRequests)

	// Get pending maintenance requests
	var pendingMaintenanceRequests int64
	config.DB.Model(&models.MaintenanceRequest{}).Where("status = ?", models.MaintenanceStatusPending).Count(&pendingMaintenanceRequests)

	// Get total reviews
	var totalReviews int64
	config.DB.Model(&models.Review{}).Count(&totalReviews)

	// Get average rating
	var avgRating float64
	config.DB.Model(&models.Review{}).Select("AVG(rating)").Scan(&avgRating)

	// Get recent activity (last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	var recentUsers int64
	config.DB.Model(&models.User{}).Where("created_at >= ?", thirtyDaysAgo).Count(&recentUsers)

	var recentHouses int64
	config.DB.Model(&models.House{}).Where("created_at >= ?", thirtyDaysAgo).Count(&recentHouses)

	var recentPayments int64
	config.DB.Model(&models.Payment{}).Where("created_at >= ?", thirtyDaysAgo).Count(&recentPayments)

	utils.SuccessResponse(c, http.StatusOK, "Dashboard statistics retrieved successfully", gin.H{
		"users": gin.H{
			"total":   totalUsers,
			"by_role": usersByRole,
			"recent":  recentUsers,
		},
		"houses": gin.H{
			"total":     totalHouses,
			"by_status": housesByStatus,
			"recent":    recentHouses,
		},
		"agreements": gin.H{
			"total":  totalAgreements,
			"active": activeAgreements,
		},
		"payments": gin.H{
			"total":      totalPayments,
			"revenue":    totalRevenue,
			"commission": totalCommission,
			"by_method":  paymentsByMethod,
			"recent":     recentPayments,
		},
		"maintenance": gin.H{
			"total":   totalMaintenanceRequests,
			"pending": pendingMaintenanceRequests,
		},
		"reviews": gin.H{
			"total":          totalReviews,
			"average_rating": avgRating,
		},
	})
}

// GetUsers handles getting all users with pagination
func (ah *AdminHandler) GetUsers(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	role := c.Query("role")
	search := c.Query("search")

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := config.DB.Model(&models.User{})

	// Apply filters
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if search != "" {
		query = query.Where("full_name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get users
	var users []models.User
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch users", err)
		return
	}

	// Remove password hashes
	for i := range users {
		users[i].PasswordHash = ""
	}

	// Calculate pagination info
	totalPages := (total + int64(limit) - 1) / int64(limit)

	utils.SuccessResponse(c, http.StatusOK, "Users retrieved successfully", gin.H{
		"users": users,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// UpdateUserStatus handles updating user status
func (ah *AdminHandler) UpdateUserStatus(c *gin.Context) {
	userID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Get user
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	// Update user status
	user.IsActive = req.IsActive
	user.UpdatedAt = time.Now()

	if err := config.DB.Save(&user).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update user status", err)
		return
	}

	user.PasswordHash = ""
	utils.SuccessResponse(c, http.StatusOK, "User status updated successfully", gin.H{
		"user": user,
	})
}

// GetReports handles getting various reports
func (ah *AdminHandler) GetReports(c *gin.Context) {
	reportType := c.Query("type")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Parse dates
	var start, end time.Time
	var err error

	if startDate != "" {
		start, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid start date format", err)
			return
		}
	} else {
		start = time.Now().AddDate(0, -1, 0) // Default to 1 month ago
	}

	if endDate != "" {
		end, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid end date format", err)
			return
		}
	} else {
		end = time.Now() // Default to now
	}

	switch reportType {
	case "payments":
		ah.getPaymentReport(c, start, end)
	case "houses":
		ah.getHouseReport(c, start, end)
	case "users":
		ah.getUserReport(c, start, end)
	default:
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid report type", nil)
	}
}

// getPaymentReport generates payment report
func (ah *AdminHandler) getPaymentReport(c *gin.Context, start, end time.Time) {
	// Get payments in date range
	var payments []models.Payment
	config.DB.Where("created_at >= ? AND created_at <= ?", start, end).
		Preload("Agreement.House").
		Preload("Agreement.Tenant").
		Find(&payments)

	// Calculate totals
	var totalAmount, totalCommission float64
	for _, payment := range payments {
		if payment.Status == models.PaymentStatusCompleted {
			totalAmount += payment.Amount
			totalCommission += payment.Commission
		}
	}

	// Get payments by method
	var paymentsByMethod []struct {
		Method string  `json:"method"`
		Count  int64   `json:"count"`
		Amount float64 `json:"amount"`
	}
	config.DB.Model(&models.Payment{}).
		Where("created_at >= ? AND created_at <= ?", start, end).
		Select("method, COUNT(*) as count, COALESCE(SUM(amount), 0) as amount").
		Group("method").
		Find(&paymentsByMethod)

	utils.SuccessResponse(c, http.StatusOK, "Payment report generated successfully", gin.H{
		"period": gin.H{
			"start_date": start.Format("2006-01-02"),
			"end_date":   end.Format("2006-01-02"),
		},
		"summary": gin.H{
			"total_amount":     totalAmount,
			"total_commission": totalCommission,
			"total_payments":   len(payments),
		},
		"payments_by_method": paymentsByMethod,
		"payments":           payments,
	})
}

// getHouseReport generates house report
func (ah *AdminHandler) getHouseReport(c *gin.Context, start, end time.Time) {
	// Get houses created in date range
	var houses []models.House
	config.DB.Where("created_at >= ? AND created_at <= ?", start, end).
		Preload("Landlord").
		Find(&houses)

	// Get houses by type
	var housesByType []struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}
	config.DB.Model(&models.House{}).
		Where("created_at >= ? AND created_at <= ?", start, end).
		Select("house_type as type, COUNT(*) as count").
		Group("house_type").
		Find(&housesByType)

	// Get houses by status
	var housesByStatus []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	config.DB.Model(&models.House{}).
		Where("created_at >= ? AND created_at <= ?", start, end).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&housesByStatus)

	utils.SuccessResponse(c, http.StatusOK, "House report generated successfully", gin.H{
		"period": gin.H{
			"start_date": start.Format("2006-01-02"),
			"end_date":   end.Format("2006-01-02"),
		},
		"summary": gin.H{
			"total_houses": len(houses),
		},
		"houses_by_type":   housesByType,
		"houses_by_status": housesByStatus,
		"houses":           houses,
	})
}

// getUserReport generates user report
func (ah *AdminHandler) getUserReport(c *gin.Context, start, end time.Time) {
	// Get users created in date range
	var users []models.User
	config.DB.Where("created_at >= ? AND created_at <= ?", start, end).Find(&users)

	// Get users by role
	var usersByRole []struct {
		Role  string `json:"role"`
		Count int64  `json:"count"`
	}
	config.DB.Model(&models.User{}).
		Where("created_at >= ? AND created_at <= ?", start, end).
		Select("role, COUNT(*) as count").
		Group("role").
		Find(&usersByRole)

	// Remove password hashes
	for i := range users {
		users[i].PasswordHash = ""
	}

	utils.SuccessResponse(c, http.StatusOK, "User report generated successfully", gin.H{
		"period": gin.H{
			"start_date": start.Format("2006-01-02"),
			"end_date":   end.Format("2006-01-02"),
		},
		"summary": gin.H{
			"total_users": len(users),
		},
		"users_by_role": usersByRole,
		"users":         users,
	})
}
