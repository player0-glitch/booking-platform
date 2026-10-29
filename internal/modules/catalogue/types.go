package catalogue

import "slices"

type status string

const (
	StatusActive   status = "active"
	StatusInActive status = "active"
)

type mode string

const (
	Duration_Variable mode = "variable"
	Duration_Fixed    mode = "fixed"
)

type pricingUnit string

const (
	Per_Booking pricingUnit = "per booking"
	Per_Unit    pricingUnit = "per unit"
	Per_Hour    pricingUnit = "per hour"
)

// godoc
// Restricts the generic functions below to types whose underlying type
// is a string
type StringEnum interface {
	~string
}

// godoc
// Generic Validation method fot enum types
func IsValidEnum[T StringEnum](value T, validOptions []T) bool {
	return slices.Contains(validOptions, value)
}
