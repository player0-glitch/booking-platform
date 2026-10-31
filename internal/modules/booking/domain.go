package booking

import (
	"booking-platform/internal/core"
	"time"

	"github.com/google/uuid"
)

type BookingAllocationStatuss struct {
	core.Base
	BookingId    uuid.UUID `gorm:"type:uuid;foreignKey"`
	ResourceId   uuid.UUID `gorm:"type:uuid;foreignKey"`
	BlockedFrom  time.Time
	BlockedUntil time.Time
	Exclusive    bool
	Quantity     int `gorm:"default:1"`
	Status       allocationStatus
}

func (BookingAllocationStatuss) TableName() string {
	return "booking_allocations"
}

type Bookings struct {
	core.Base
	OrganisationId  uuid.UUID `gorm:"type:uuid;foreignKey"`
	CustomerId      uuid.UUID `gorm:"type:uuid;foreignKey"`
	ServicedId      uuid.UUID `gorm:"type:uuid;foreignKey"`
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
