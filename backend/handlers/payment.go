package handlers

import (
	"bondihub/config"
	"bondihub/models"
	"bondihub/services"
	"bondihub/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PaymentHandler handles payment-related requests
type PaymentHandler struct {
	paymentService *services.PaymentService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{
		paymentService: services.NewPaymentService(),
	}
}

// CreatePaymentRequest represents the request structure for creating a payment
type CreatePaymentRequest struct {
	AgreementID uuid.UUID `json:"agreement_id" binding:"required"`
	Amount      float64   `json:"amount" binding:"required,min=0"`
	Method      string    `json:"method" binding:"required,oneof=MTN Airtel Cash Bank"`
	ReferenceNo string    `json:"reference_no"`
}

// ProcessPayment handles processing a payment
// ProcessPayment processes a payment for a rental agreement
// @Summary Process payment
// @Description Process a payment for a rental agreement using various payment methods
// @Tags Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePaymentRequest true "Payment details"
// @Success 201 {object} map[string]interface{} "Payment processed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /payments [post]
func (ph *PaymentHandler) ProcessPayment(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)

	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Get rental agreement
	var agreement models.RentalAgreement
	if err := config.DB.Preload("House").Preload("Tenant").First(&agreement, req.AgreementID).Error; err != nil {
		utils.NotFoundResponse(c, "Rental agreement not found")
		return
	}

	// Check if user is the tenant or admin
	if agreement.TenantID != userModel.ID && userModel.Role != models.RoleAdmin {
		utils.ForbiddenResponse(c, "You can only make payments for your own agreements")
		return
	}

	// Check if agreement is active
	if agreement.Status != models.AgreementStatusActive {
		utils.ErrorResponse(c, http.StatusBadRequest, "Cannot make payment for inactive agreement", nil)
		return
	}

	// Generate reference number if not provided
	if req.ReferenceNo == "" {
		req.ReferenceNo = fmt.Sprintf("PAY_%d", time.Now().Unix())
	}

	// Create payment record
	payment := models.Payment{
		AgreementID: agreement.ID,
		Amount:      req.Amount,
		PaymentDate: time.Now(),
		Method:      models.PaymentMethod(req.Method),
		ReferenceNo: req.ReferenceNo,
		Status:      models.PaymentStatusPending,
		Commission:  ph.paymentService.CalculateCommission(req.Amount),
	}

	if err := config.DB.Create(&payment).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create payment record", err)
		return
	}

	// Process payment
	result, err := ph.paymentService.ProcessPayment(&payment, payment.Method)
	if err != nil {
		// Update payment status to failed
		payment.Status = models.PaymentStatusFailed
		config.DB.Save(&payment)

		utils.InternalServerErrorResponse(c, "Payment processing failed", err)
		return
	}

	// Update payment status based on result
	if result.Success {
		payment.Status = models.PaymentStatusCompleted
		payment.ReferenceNo = result.ReferenceNo
	} else {
		payment.Status = models.PaymentStatusFailed
	}

	config.DB.Save(&payment)

	// Create notification for landlord
	notification := models.Notification{
		UserID:  agreement.House.LandlordID,
		Title:   "New Payment Received",
		Message: fmt.Sprintf("Payment of ZMW %.2f received for %s", payment.Amount, agreement.House.Title),
		Type:    "payment",
	}
	config.DB.Create(&notification)

	utils.SuccessResponse(c, http.StatusCreated, "Payment processed successfully", gin.H{
		"payment": payment,
		"result":  result,
	})
}

// GetPayments handles getting payments for a user
// GetPayments retrieves user's payment history
// @Summary Get payments
// @Description Get payment history for the authenticated user
// @Tags Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{} "Payments retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /payments [get]
func (ph *PaymentHandler) GetPayments(c *gin.Context) {
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
	method := c.Query("method")

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := config.DB.Model(&models.Payment{}).
		Joins("JOIN rental_agreements ON payments.agreement_id = rental_agreements.id").
		Preload("Agreement.House").
		Preload("Agreement.Tenant")

	// Apply filters based on user role
	if userModel.Role == models.RoleTenant {
		query = query.Where("rental_agreements.tenant_id = ?", userModel.ID)
	} else if userModel.Role == models.RoleLandlord {
		query = query.Joins("JOIN houses ON rental_agreements.house_id = houses.id").
			Where("houses.landlord_id = ?", userModel.ID)
	}

	if status != "" {
		query = query.Where("payments.status = ?", status)
	}
	if method != "" {
		query = query.Where("payments.method = ?", method)
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get payments
	var payments []models.Payment
	if err := query.Offset(offset).Limit(limit).Order("payments.created_at DESC").Find(&payments).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch payments", err)
		return
	}

	// Calculate pagination info
	totalPages := (total + int64(limit) - 1) / int64(limit)

	utils.SuccessResponse(c, http.StatusOK, "Payments retrieved successfully", gin.H{
		"payments": payments,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetPayment handles getting a single payment by ID
func (ph *PaymentHandler) GetPayment(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	paymentID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(paymentID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payment ID", err)
		return
	}

	var payment models.Payment
	if err := config.DB.Preload("Agreement.House").Preload("Agreement.Tenant").First(&payment, id).Error; err != nil {
		utils.NotFoundResponse(c, "Payment not found")
		return
	}

	// Check if user has access to this payment
	hasAccess := false
	if userModel.Role == models.RoleAdmin {
		hasAccess = true
	} else if userModel.Role == models.RoleTenant && payment.Agreement.TenantID == userModel.ID {
		hasAccess = true
	} else if userModel.Role == models.RoleLandlord && payment.Agreement.House.LandlordID == userModel.ID {
		hasAccess = true
	}

	if !hasAccess {
		utils.ForbiddenResponse(c, "You don't have access to this payment")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment retrieved successfully", gin.H{
		"payment": payment,
	})
}

// GetPaymentStats handles getting payment statistics
func (ph *PaymentHandler) GetPaymentStats(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)

	// Build base query
	query := config.DB.Model(&models.Payment{}).
		Joins("JOIN rental_agreements ON payments.agreement_id = rental_agreements.id")

	// Apply filters based on user role
	if userModel.Role == models.RoleTenant {
		query = query.Where("rental_agreements.tenant_id = ?", userModel.ID)
	} else if userModel.Role == models.RoleLandlord {
		query = query.Joins("JOIN houses ON rental_agreements.house_id = houses.id").
			Where("houses.landlord_id = ?", userModel.ID)
	}

	// Get total payments
	var totalPayments int64
	query.Count(&totalPayments)

	// Get total amount
	var totalAmount float64
	query.Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)

	// Get completed payments
	var completedPayments int64
	query.Where("status = ?", models.PaymentStatusCompleted).Count(&completedPayments)

	// Get completed amount
	var completedAmount float64
	query.Where("status = ?", models.PaymentStatusCompleted).Select("COALESCE(SUM(amount), 0)").Scan(&completedAmount)

	// Get pending payments
	var pendingPayments int64
	query.Where("status = ?", models.PaymentStatusPending).Count(&pendingPayments)

	// Get failed payments
	var failedPayments int64
	query.Where("status = ?", models.PaymentStatusFailed).Count(&failedPayments)

	// Get payments by method
	var paymentsByMethod []struct {
		Method string  `json:"method"`
		Count  int64   `json:"count"`
		Amount float64 `json:"amount"`
	}
	query.Select("method, COUNT(*) as count, COALESCE(SUM(amount), 0) as amount").
		Group("method").
		Find(&paymentsByMethod)

	utils.SuccessResponse(c, http.StatusOK, "Payment statistics retrieved successfully", gin.H{
		"total_payments":     totalPayments,
		"total_amount":       totalAmount,
		"completed_payments": completedPayments,
		"completed_amount":   completedAmount,
		"pending_payments":   pendingPayments,
		"failed_payments":    failedPayments,
		"payments_by_method": paymentsByMethod,
	})
}
