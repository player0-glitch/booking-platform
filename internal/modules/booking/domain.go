// Package booking
package booking

import (
	"time"

	"booking-platform/internal/core"

	"github.com/google/uuid"
)

type allocationStatus string

const (
	allocationStatusHeld     allocationStatus = "held"
	allocationStatusActive   allocationStatus = "active"
	allocationStatusReleased allocationStatus = "released"
)

type BookingAllocationStatuss struct {
	core.Base
	BookingID    uuid.UUID `gorm:"type:uuid;foreignKey"`
	ResourceID   uuid.UUID `gorm:"type:uuid;foreignKey"`
	BlockedFrom  time.Time
	BlockedUntil time.Time
	Exclusive    bool
	Quantity     int `gorm:"default:1"`
	Status       allocationStatus
}

func (BookingAllocationStatuss) TableName() string {
	return "booking_allocations"
}

type bookingStatus string

const (
	bookingPending   bookingStatus = "pending"
	bookingConfirmed bookingStatus = "confirmed"
	bookingDeclined  bookingStatus = "declined"
	bookingCancelled bookingStatus = "cancelled"
	bookingCompleted bookingStatus = "completed"
	bookingNoShow    bookingStatus = "no show"
	bookingExpired   bookingStatus = "expired"
)

type Bookings struct {
	core.Base
	OrganisationID  uuid.UUID `gorm:"type:uuid;foreignKey"`
	CustomerID      uuid.UUID `gorm:"type:uuid;foreignKey"`
	ServicedID      uuid.UUID `gorm:"type:uuid;foreignKey"`
	StartAt         time.Time
	EndsAt          time.Time
	Notes           string
	DurationMinutes time.Duration
	TotalPrice      int    `gorm:"type:int"`
	Currency        string `gorm:"type:char(3)"`
	ExpiresAt       *time.Time
	PartySize       int `gorm:"default:1"`
	Status          bookingStatus
}

func (Bookings) TableName() string {
	return "bookings"
}
