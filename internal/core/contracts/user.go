package core

import (
	"context"
)

type User struct {
	Id    int
	Email string
}

// This is what the booking module will use to communcate with the persistence of user
type UserReader interface {
	FindById(context.Context, uint) (*User, error)
}
