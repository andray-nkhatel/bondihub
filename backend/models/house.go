package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HouseStatus represents the status of a house
type HouseStatus string

const (
	StatusAvailable   HouseStatus = "available"
	StatusOccupied    HouseStatus = "occupied"
	StatusMaintenance HouseStatus = "maintenance"
)

// HouseType represents the type of house
type HouseType string

const (
	TypeApartment  HouseType = "apartment"
	TypeHouse      HouseType = "house"
	TypeStudio     HouseType = "studio"
	TypeTownhouse  HouseType = "townhouse"
	TypeCommercial HouseType = "commercial"
)

// House represents a property listing
type House struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LandlordID    uuid.UUID      `json:"landlord_id" gorm:"type:uuid;not null"`
	Title         string         `json:"title" gorm:"not null"`
	Description   string         `json:"description" gorm:"type:text"`
	Address       string         `json:"address" gorm:"not null"`
	MonthlyRent   float64        `json:"monthly_rent" gorm:"not null;type:decimal(10,2)"`
	Status        HouseStatus    `json:"status" gorm:"not null;default:'available'"`
	HouseType     HouseType      `json:"house_type" gorm:"not null"`
	Latitude      float64        `json:"latitude" gorm:"type:decimal(10,8)"`
	Longitude     float64        `json:"longitude" gorm:"type:decimal(11,8)"`
	Bedrooms      int            `json:"bedrooms" gorm:"default:0"`
	Bathrooms     int            `json:"bathrooms" gorm:"default:0"`
	Area          float64        `json:"area" gorm:"type:decimal(8,2)"` // in square meters
	IsFeatured    bool           `json:"is_featured" gorm:"default:false"`
	FeaturedUntil *time.Time     `json:"featured_until"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Landlord            User                 `json:"landlord,omitempty" gorm:"foreignKey:LandlordID"`
	Images              []HouseImage         `json:"images,omitempty" gorm:"foreignKey:HouseID"`
	RentalAgreements    []RentalAgreement    `json:"rental_agreements,omitempty" gorm:"foreignKey:HouseID"`
	Reviews             []Review             `json:"reviews,omitempty" gorm:"foreignKey:HouseID"`
	MaintenanceRequests []MaintenanceRequest `json:"maintenance_requests,omitempty" gorm:"foreignKey:HouseID"`
	Favorites           []Favorite           `json:"favorites,omitempty" gorm:"foreignKey:HouseID"`
}

// BeforeCreate hook to set default values
func (h *House) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for House
func (House) TableName() string {
	return "houses"
}

// HouseImage represents images associated with a house
type HouseImage struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	HouseID   uuid.UUID `json:"house_id" gorm:"type:uuid;not null"`
	ImageURL  string    `json:"image_url" gorm:"not null"`
	IsPrimary bool      `json:"is_primary" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	House House `json:"house,omitempty" gorm:"foreignKey:HouseID"`
}

// BeforeCreate hook to set default values
func (hi *HouseImage) BeforeCreate(tx *gorm.DB) error {
	if hi.ID == uuid.Nil {
		hi.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for HouseImage
func (HouseImage) TableName() string {
	return "house_images"
}
