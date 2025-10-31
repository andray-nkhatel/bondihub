package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Review represents a tenant's review of a house
type Review struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID  uuid.UUID      `json:"tenant_id" gorm:"type:uuid;not null"`
	HouseID   uuid.UUID      `json:"house_id" gorm:"type:uuid;not null"`
	Rating    int            `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Comment   string         `json:"comment" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Tenant User  `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
	House  House `json:"house,omitempty" gorm:"foreignKey:HouseID"`
}

// BeforeCreate hook to set default values
func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for Review
func (Review) TableName() string {
	return "reviews"
}

// MaintenanceRequestStatus represents the status of a maintenance request
type MaintenanceRequestStatus string

const (
	MaintenanceStatusPending    MaintenanceRequestStatus = "pending"
	MaintenanceStatusInProgress MaintenanceRequestStatus = "in_progress"
	MaintenanceStatusResolved   MaintenanceRequestStatus = "resolved"
	MaintenanceStatusCancelled  MaintenanceRequestStatus = "cancelled"
)

// MaintenanceRequest represents a maintenance request from a tenant
type MaintenanceRequest struct {
	ID          uuid.UUID                `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID                `json:"tenant_id" gorm:"type:uuid;not null"`
	HouseID     uuid.UUID                `json:"house_id" gorm:"type:uuid;not null"`
	Title       string                   `json:"title" gorm:"not null"`
	Description string                   `json:"description" gorm:"type:text;not null"`
	Status      MaintenanceRequestStatus `json:"status" gorm:"not null;default:'pending'"`
	Priority    string                   `json:"priority" gorm:"default:'medium'"` // low, medium, high, urgent
	ReportedAt  time.Time                `json:"reported_at" gorm:"not null"`
	ResolvedAt  *time.Time               `json:"resolved_at"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
	DeletedAt   gorm.DeletedAt           `json:"-" gorm:"index"`

	// Relationships
	Tenant User  `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
	House  House `json:"house,omitempty" gorm:"foreignKey:HouseID"`
}

// BeforeCreate hook to set default values
func (mr *MaintenanceRequest) BeforeCreate(tx *gorm.DB) error {
	if mr.ID == uuid.Nil {
		mr.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for MaintenanceRequest
func (MaintenanceRequest) TableName() string {
	return "maintenance_requests"
}

// Favorite represents a tenant's favorite house
type Favorite struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID  uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null"`
	HouseID   uuid.UUID `json:"house_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Tenant User  `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
	House  House `json:"house,omitempty" gorm:"foreignKey:HouseID"`
}

// BeforeCreate hook to set default values
func (f *Favorite) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for Favorite
func (Favorite) TableName() string {
	return "favorites"
}

// Notification represents a notification for a user
type Notification struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Title     string    `json:"title" gorm:"not null"`
	Message   string    `json:"message" gorm:"type:text;not null"`
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	Type      string    `json:"type" gorm:"not null"` // payment, maintenance, agreement, general
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// BeforeCreate hook to set default values
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for Notification
func (Notification) TableName() string {
	return "notifications"
}
