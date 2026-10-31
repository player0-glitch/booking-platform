package booking

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

type allocationStatus string

const (
	allocationStatusHeld     allocationStatus = "held"
	allocationStatusActive   allocationStatus = "active"
	allocationStatusReleased allocationStatus = "released"
)
