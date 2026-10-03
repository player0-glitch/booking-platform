package auth

import (
	"context"
	"net/http"
)

// Private to the auth package
//
//	type authenticator interface {
//		Authenticate(email, password string) (string, error)
//	}
type storeManager interface {
	SaveToken(w http.ResponseWriter, r *http.Request, token string) error
	Clear(w http.ResponseWriter, r *http.Request) error
}

type UserRole struct {
	Id   uint
	Name string
}

// User struct with the attributes expected byt the Auth Domain
// This is adapted to from the User model
type AuthUser struct {
	Id           uint
	Email        string
	PasswordHash string
	Roles        []UserRole
}

// Contract that the Auth Domain expect to read from the User Domain
type UserReader interface {
	GetByEmailWithRoles(context.Context, string) (*AuthUser, error)
}
