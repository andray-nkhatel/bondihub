package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole represents the different user roles in the system
type UserRole string

const (
	RoleLandlord UserRole = "landlord"
	RoleTenant   UserRole = "tenant"
	RoleAgent    UserRole = "agent"
	RoleAdmin    UserRole = "admin"
)

// SubscriptionPlan represents subscription plans for landlords/agents
type SubscriptionPlan string

const (
	PlanBasic      SubscriptionPlan = "basic"
	PlanPremium    SubscriptionPlan = "premium"
	PlanEnterprise SubscriptionPlan = "enterprise"
)

// User represents a user in the system
type User struct {
	ID               uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FullName         string           `json:"full_name" gorm:"not null"`
	Email            string           `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash     string           `json:"-" gorm:"not null"`
	Phone            string           `json:"phone" gorm:"not null"`
	Role             UserRole         `json:"role" gorm:"not null;default:'tenant'"`
	IsActive         bool             `json:"is_active" gorm:"default:true"`
	IsVerified       bool             `json:"is_verified" gorm:"default:false"`
	ProfileImage     string           `json:"profile_image"`
	SubscriptionPlan SubscriptionPlan `json:"subscription_plan" gorm:"default:'basic'"`
	PlanExpiryDate   *time.Time       `json:"plan_expiry_date"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	DeletedAt        gorm.DeletedAt   `json:"-" gorm:"index"`

	// Relationships
	HousesAsLandlord    []House              `json:"houses,omitempty" gorm:"foreignKey:LandlordID"`
	RentalAgreements    []RentalAgreement    `json:"rental_agreements,omitempty" gorm:"foreignKey:TenantID"`
	Reviews             []Review             `json:"reviews,omitempty" gorm:"foreignKey:TenantID"`
	MaintenanceRequests []MaintenanceRequest `json:"maintenance_requests,omitempty" gorm:"foreignKey:TenantID"`
	Favorites           []Favorite           `json:"favorites,omitempty" gorm:"foreignKey:TenantID"`
	Notifications       []Notification       `json:"notifications,omitempty" gorm:"foreignKey:UserID"`
}

// BeforeCreate hook to set default values
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}
