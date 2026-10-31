package catalogue

type status string

const (
	statusActive   status = "active"
	statusInActive status = "active"
)

type mode string

const (
	duration_Variable mode = "variable"
	duration_Fixed    mode = "fixed"
)

type pricingUnit string

const (
	per_Booking pricingUnit = "per booking"
	per_Unit    pricingUnit = "per unit"
	per_Hour    pricingUnit = "per hour"
)
