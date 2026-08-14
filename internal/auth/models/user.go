// auth/models/user.go
package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName    string         `gorm:"not null"`
	LastName     string         `gorm:"not null"`
	Email        string         `gorm:"uniqueIndex;not null"`
	PasswordHash string         `gorm:"not null"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
