// Package catalogue orgranisation,resource,services make up the business catalogue
package catalogue

import (
	"booking-platform/internal/core"

	"github.com/google/uuid"
)

type status string

const (
	statusActive   status = "active"
	statusInActive status = "inactive"
)

type Organisation struct {
	Base     core.Base
	Slug     string `gorm:"uniqueIndex;not null"`
	Status   status `gorm:"type:text;default:'active'"`
	TimeZone string
}

func (Organisation) TableName() string {
	return "organisations"
}

type Resource struct {
	Base           core.Base
	OrganisationID uuid.UUID `gorm:"uniqueIndex:idx_organisation_id_name"`
	Kind           string
	Name           string `gorm:"uniqueIndex:idx_organisation_id_name"`
	BookingMode    string `gorm:"default:exclusive"`
	Capacity       int
	Status         status `gorm:"type:text;default:'active'"`
	// This is for eager loading in joining query
	Organisation Organisation `gorm:"foreignKey:OrganisationID"`
}

func (Resource) TableName() string {
	return "resources"
}

type mode string

const (
	durationVariable mode = "variable"
	durationFixed    mode = "fixed"
)

type pricingUnit string

const (
	per_Booking pricingUnit = "per booking"
	per_Unit    pricingUnit = "per unit"
	per_Hour    pricingUnit = "per hour"
)

type Service struct {
	Base                core.Base
	OrganisationID      uuid.UUID `gorm:"uniqueIndex:idx_organisation_id_name"`
	Name                string    `gorm:"uniqueIndex:idx_organisation_id_name"`
	DurationMode        mode      `gorm:"type:text"`
	DurationMinutes     int
	MinDurationMinutes  *int
	MaxDurationMinutes  *int
	DurationStep        *int
	BufferBeforeMinutes int
	BufferAfterMinutes  int
	PriceMinor          int
	PricingUnit         pricingUnit `gorm:"type:text"`
	Currency            string      `gorm:"type:char(3)"`
	RequiresApproval    bool        `gorm:"default:false"`
	Status              status      `gorm:"type:text;default:'active'"`
	// This is for eager loading in joining query
	Organisation Organisation `gorm:"foreignKey:OrganisationID"`
}

func (Service) TableName() string {
	return "services"
}

type ResourceServices struct {
	// Composite Keys
	ResourceID         uuid.UUID `gorm:"type:uuid;primaryKey;column:resource_id"`
	ServiceID          uuid.UUID `gorm:"type:uuid;primaryKey;index:idx_resource_services_service_id;column:service_id"`
	PriceOverrideMinor *int
}

func (ResourceServices) TableName() string {
	return "resource_services"
}
