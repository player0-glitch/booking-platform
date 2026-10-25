package user

import (
	"booking-platform/internal/core"
	"time"

	"gorm.io/gorm"
)

type User struct {
	core.Base
	FirstName    string `gorm:"not null"`
	LastName     string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	//This is for earger loading in a joining query
	RoleId int  `gorm:"not null;default:3"`
	Role   Role `gorm:"foreignKey:RoleId"`
}
type Role struct {
	Id        int
	CreatedAt time.Time
	UpdatedAt time.Time
	RoleName  string `gorm:"unique;not null" json:"role_name"`

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
