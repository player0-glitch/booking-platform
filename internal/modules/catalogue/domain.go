package catalogue

import (
	"booking-platform/internal/core"

	"github.com/google/uuid"
)

// orgranisation,resource,services make up the business catalogue

type Orgranisation struct {
	Base     core.Base
	Slug     string `gorm:"uniqueIndex;not null"`
	Status   status `gorm:"type:text;default:'active'"`
	TimeZone string
}

// Tells gorm which sql table to map from
func (Orgranisation) TableName() string {
	return "organisations"
}

type Resource struct {
	Base            core.Base
	OrgranisationId uuid.UUID `gorm:"uniqueIndex:idx_organisation_id_name"`
	Kind            string
	Name            string `gorm:"uniqueIndex:idx_organisation_id_name"`
	BookingMode     string `gorm:"default:exclusive"`
	Capacity        int
	Status          status `gorm:"type:text;default:'active'"`
	//This is for eager loading in joining query
	Orgranisation Orgranisation `gorm:"foreignKey:OrganisationId"`
}

func (Resource) TableName() string {
	return "resources"
}

type Service struct {
	Base                core.Base
	OrgranisationId     uuid.UUID `gorm:"uniqueIndex:idx_organisation_id_name"`
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
	//This is for eager loading in joining query
	Orgranisation Orgranisation `gorm:"foreignKey:OrganisationId"`
}

func (Service) TableName() string {
	return "services"
}

type ResourceServices struct {
	//Composite Keys
	ResourceId         uuid.UUID `gorm:"type:uuid;primaryKey;column:resource_id"`
	ServiceId          uuid.UUID `gorm:"type:uuid;primaryKey;index:idx_resource_services_service_id;column:service_id"`
	PriceOverrideMinor *int
}

func (ResourceServices) TableName() string {
	return "resource_services"
}
