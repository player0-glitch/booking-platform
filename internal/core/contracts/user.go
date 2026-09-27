package core

import (
	"context"
)

type Role struct {
	Name string
}
type User struct {
	Id           int
	Email        string
	PasswordHash string
	Roles        []Role
}

// This is what the booking module will use to communcate with the persistence of user
type UserReader interface {
	FindById(context.Context, uint) (*User, error)
	GetByEmail(context.Context, string) (*User, error)
	GetByEmailWithRoles(context.Context, string) (*User, error)
}
