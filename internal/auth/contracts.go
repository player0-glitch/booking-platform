package auth

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type storeManager interface {
	SaveToken(w http.ResponseWriter, r *http.Request, token string) error
	Clear(w http.ResponseWriter, r *http.Request) error
}

type UserRole struct {
	Id   int
	Name string
}

// User struct with the attributes expected byt the Auth Domain
// This is adapted to from the User model
type AuthUser struct {
	Id           uuid.UUID
	Email        string
	PasswordHash string
	Role         UserRole
}

// Contract that the Auth Domain expect to read from the User Domain
type UserReader interface {
	GetByEmailWithRoles(context.Context, string) (*AuthUser, error)
	CreateProvisonGuest()
}
