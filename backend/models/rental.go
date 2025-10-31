package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AgreementStatus represents the status of a rental agreement
type AgreementStatus string

const (
	AgreementStatusActive     AgreementStatus = "active"
	AgreementStatusTerminated AgreementStatus = "terminated"
	AgreementStatusExpired    AgreementStatus = "expired"
)

// RentalAgreement represents a rental agreement between landlord and tenant
type RentalAgreement struct {
	ID         uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	HouseID    uuid.UUID       `json:"house_id" gorm:"type:uuid;not null"`
	TenantID   uuid.UUID       `json:"tenant_id" gorm:"type:uuid;not null"`
	StartDate  time.Time       `json:"start_date" gorm:"not null"`
	EndDate    time.Time       `json:"end_date" gorm:"not null"`
	RentAmount float64         `json:"rent_amount" gorm:"not null;type:decimal(10,2)"`
	Deposit    float64         `json:"deposit" gorm:"not null;type:decimal(10,2)"`
	Status     AgreementStatus `json:"status" gorm:"not null;default:'active'"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	DeletedAt  gorm.DeletedAt  `json:"-" gorm:"index"`

	// Relationships
	House    House     `json:"house,omitempty" gorm:"foreignKey:HouseID"`
	Tenant   User      `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
	Payments []Payment `json:"payments,omitempty" gorm:"foreignKey:AgreementID"`
}

// BeforeCreate hook to set default values
func (ra *RentalAgreement) BeforeCreate(tx *gorm.DB) error {
	if ra.ID == uuid.Nil {
		ra.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for RentalAgreement
func (RentalAgreement) TableName() string {
	return "rental_agreements"
}

// PaymentMethod represents the payment method used
type PaymentMethod string

const (
	PaymentMethodMTN    PaymentMethod = "MTN"
	PaymentMethodAirtel PaymentMethod = "Airtel"
	PaymentMethodCash   PaymentMethod = "Cash"
	PaymentMethodBank   PaymentMethod = "Bank"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// Payment represents a rent payment
type Payment struct {
	ID          uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AgreementID uuid.UUID     `json:"agreement_id" gorm:"type:uuid;not null"`
	Amount      float64       `json:"amount" gorm:"not null;type:decimal(10,2)"`
	PaymentDate time.Time     `json:"payment_date" gorm:"not null"`
	Method      PaymentMethod `json:"method" gorm:"not null"`
	ReferenceNo string        `json:"reference_no" gorm:"uniqueIndex"`
	Status      PaymentStatus `json:"status" gorm:"not null;default:'pending'"`
	Commission  float64       `json:"commission" gorm:"type:decimal(10,2);default:0"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`

	// Relationships
	Agreement RentalAgreement `json:"agreement,omitempty" gorm:"foreignKey:AgreementID"`
}

// BeforeCreate hook to set default values
func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for Payment
func (Payment) TableName() string {
	return "payments"
}
