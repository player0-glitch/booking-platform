package user

import (
	"booking-platform/internal/core"

	"gorm.io/gorm"
)

type AccountStatus string

const (
	AccountGuest      AccountStatus = "guest"
	AccountRegistered AccountStatus = "registered"
	AccountSuspended  AccountStatus = "suspended"
)

// godoc
// A User is an identity on the platform,'who is this'. A guest is an anonymous user (imagine
// A user DOES not have a role. Nullable options accommodate guest user
// a potential customer walk in but they don't spend money)
type User struct {
	core.Base
	FirstName     *string `gorm:"type:varchar(255)"`
	LastName      *string `gorm:"type:varchar(255)"`
	Email         *string `gorm:"type:varchar(255)"`
	PasswordHash  *string
	AccountStatus AccountStatus `gorm:"default:'registered'"`
	//This is for earger loading in a joining query
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
