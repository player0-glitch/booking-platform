package user

import "gorm.io/gorm"

type User struct {
	gorm.Model          //id,created_at,updated_at
	FirstName    string `gorm:"not null"`
	LastName     string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	//This is for earger loading in a joining query
	Roles     []Role         `gorm:"many2many;user_roles"`
	DeletedAt gorm.DeletedAt `gorm:"index;null"`
}

type Role struct {
	gorm.Model        //makes id,created_at,updated_at and deleted_at
	RoleName   string `gorm:"unique;not null" json:"role_name"`
}
