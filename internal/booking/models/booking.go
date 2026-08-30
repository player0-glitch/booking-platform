// models/booking.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model //creates ID,Created_At,Updated_At
	Start      *time.Time
	End        *time.Time
	// Guest     Guest
	DeletedAt gorm.DeletedAt `gorm:"index"` //can be used to index using this date
}
