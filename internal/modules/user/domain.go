package user

import (
	"booking-platform/internal/core"
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	// gorm.Model          //id,created_at,updated_at
	core.Base
	FirstName    string `gorm:"not null"`
	LastName     string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	//This is for earger loading in a joining query
	Roles     []Role         `gorm:"many2many:user_roles"`
	DeletedAt gorm.DeletedAt `gorm:"index;null"`
}

/*
  - Perhaps only repository with a 'Parent' service layer should not
  - be given an interface and rather that have a concrete struct implementation.
  - User Interface that an unrelated service needs or as a domain contract
    type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindById(ctx context.Context, id uint) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    GetByEmailWithRoles(ctx context.Context, email string) (*User, error)
    FindAll(ctx context.Context) ([]User, error)
    DeleteByEmail(ctx context.Context, email string) error
    DeleteById(ctx context.Context, id uint) error
    SoftDeleteById(ctx context.Context, id uint) error
    }
*/
type Role struct {
	Id        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	// core.Base
	RoleName string `gorm:"unique;not null" json:"role_name"`
}

type RolesRepository interface {
	RemoveRoleFromUser(ctx context.Context, user *User, role *Role) error
	ChangeUserRole(ctx context.Context, user *User, role *Role) error
}
