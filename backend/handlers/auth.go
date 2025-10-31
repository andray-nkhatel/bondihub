package handlers

import (
	"bondihub/config"
	"bondihub/models"
	"bondihub/services"
	"bondihub/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	paymentService *services.PaymentService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		paymentService: services.NewPaymentService(),
	}
}

// RegisterRequest represents the request structure for user registration
type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=100" example:"John Doe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
	Phone    string `json:"phone" binding:"required,min=10" example:"+260123456789"`
	Role     string `json:"role" binding:"required,oneof=landlord tenant agent admin" example:"tenant"`
}

// LoginRequest represents the request structure for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration data"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 409 {object} map[string]interface{} "User already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/register [post]
func (ah *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "User with this email already exists", nil)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to process password", err)
		return
	}

	// Create user
	user := models.User{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Phone:        req.Phone,
		Role:         models.UserRole(req.Role),
		IsActive:     true,
		IsVerified:   false,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create user", err)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate token", err)
		return
	}

	// Return user data without password
	user.PasswordHash = ""
	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", gin.H{
		"user":  user,
		"token": token,
	})
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login [post]
func (ah *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Find user by email
	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.UnauthorizedResponse(c, "Invalid email or password")
		return
	}

	// Check if user is active
	if !user.IsActive {
		utils.UnauthorizedResponse(c, "Account is deactivated")
		return
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		utils.UnauthorizedResponse(c, "Invalid email or password")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate token", err)
		return
	}

	// Update last login (you might want to add this field to User model)
	user.UpdatedAt = time.Now()
	config.DB.Save(&user)

	// Return user data without password
	user.PasswordHash = ""
	utils.SuccessResponse(c, http.StatusOK, "Login successful", gin.H{
		"user":  user,
		"token": token,
	})
}

// GetProfile handles getting user profile
func (ah *AuthHandler) GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	userModel.PasswordHash = ""

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", gin.H{
		"user": userModel,
	})
}

// UpdateProfile handles updating user profile
func (ah *AuthHandler) UpdateProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)

	var req struct {
		FullName     string `json:"full_name"`
		Phone        string `json:"phone"`
		ProfileImage string `json:"profile_image"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Update user fields
	if req.FullName != "" {
		userModel.FullName = req.FullName
	}
	if req.Phone != "" {
		userModel.Phone = req.Phone
	}
	if req.ProfileImage != "" {
		userModel.ProfileImage = req.ProfileImage
	}

	userModel.UpdatedAt = time.Now()

	if err := config.DB.Save(&userModel).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update profile", err)
		return
	}

	userModel.PasswordHash = ""
	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", gin.H{
		"user": userModel,
	})
}

// ChangePassword handles changing user password
func (ah *AuthHandler) ChangePassword(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Verify current password
	if !utils.CheckPasswordHash(req.CurrentPassword, userModel.PasswordHash) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Current password is incorrect", nil)
		return
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to process new password", err)
		return
	}

	// Update password
	userModel.PasswordHash = hashedPassword
	userModel.UpdatedAt = time.Now()

	if err := config.DB.Save(&userModel).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update password", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password changed successfully", nil)
}

// Logout handles user logout (client-side token removal)
func (ah *AuthHandler) Logout(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Logout successful", nil)
}
