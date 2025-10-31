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

// NotificationHandler handles notification-related requests
type NotificationHandler struct{}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

// GetNotifications handles getting user notifications
func (nh *NotificationHandler) GetNotifications(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	unreadOnly := c.Query("unread_only") == "true"
	notificationType := c.Query("type")

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := config.DB.Model(&models.Notification{}).
		Where("user_id = ?", userModel.ID)

	// Apply filters
	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}
	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get notifications
	var notifications []models.Notification
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&notifications).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch notifications", err)
		return
	}

	// Calculate pagination info
	totalPages := (total + int64(limit) - 1) / int64(limit)

	utils.SuccessResponse(c, http.StatusOK, "Notifications retrieved successfully", gin.H{
		"notifications": notifications,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetNotification handles getting a single notification by ID
func (nh *NotificationHandler) GetNotification(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	notificationID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(notificationID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid notification ID", err)
		return
	}

	var notification models.Notification
	if err := config.DB.Where("id = ? AND user_id = ?", id, userModel.ID).First(&notification).Error; err != nil {
		utils.NotFoundResponse(c, "Notification not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notification retrieved successfully", gin.H{
		"notification": notification,
	})
}

// MarkAsRead handles marking a notification as read
func (nh *NotificationHandler) MarkAsRead(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	notificationID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(notificationID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid notification ID", err)
		return
	}

	var notification models.Notification
	if err := config.DB.Where("id = ? AND user_id = ?", id, userModel.ID).First(&notification).Error; err != nil {
		utils.NotFoundResponse(c, "Notification not found")
		return
	}

	// Mark as read
	notification.IsRead = true
	if err := config.DB.Save(&notification).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to mark notification as read", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notification marked as read successfully", gin.H{
		"notification": notification,
	})
}

// MarkAllAsRead handles marking all notifications as read
func (nh *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)

	// Mark all notifications as read
	if err := config.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userModel.ID, false).
		Update("is_read", true).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to mark all notifications as read", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "All notifications marked as read successfully", nil)
}

// DeleteNotification handles deleting a notification
func (nh *NotificationHandler) DeleteNotification(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)
	notificationID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(notificationID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid notification ID", err)
		return
	}

	var notification models.Notification
	if err := config.DB.Where("id = ? AND user_id = ?", id, userModel.ID).First(&notification).Error; err != nil {
		utils.NotFoundResponse(c, "Notification not found")
		return
	}

	// Delete notification
	if err := config.DB.Delete(&notification).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete notification", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notification deleted successfully", nil)
}

// GetNotificationStats handles getting notification statistics
func (nh *NotificationHandler) GetNotificationStats(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userModel := user.(models.User)

	// Get total notifications
	var totalNotifications int64
	config.DB.Model(&models.Notification{}).Where("user_id = ?", userModel.ID).Count(&totalNotifications)

	// Get unread notifications
	var unreadNotifications int64
	config.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userModel.ID, false).Count(&unreadNotifications)

	// Get notifications by type
	var notificationsByType []struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}
	config.DB.Model(&models.Notification{}).
		Where("user_id = ?", userModel.ID).
		Select("type, COUNT(*) as count").
		Group("type").
		Find(&notificationsByType)

	utils.SuccessResponse(c, http.StatusOK, "Notification statistics retrieved successfully", gin.H{
		"total_notifications":   totalNotifications,
		"unread_notifications":  unreadNotifications,
		"notifications_by_type": notificationsByType,
	})
}
