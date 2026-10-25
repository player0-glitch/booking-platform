package auth

import (
	"booking-platform/internal/core"
	"context"
	"net/http"

	"github.com/google/uuid"
)

type storeManagerInterface interface {
	SaveToken(w http.ResponseWriter, r *http.Request, token string) error
	Clear(w http.ResponseWriter, r *http.Request) error
}

// User struct with the attributes expected byt the Auth Domain
// This is adapted to from the User model
type AuthUser struct {
	Id           uuid.UUID
	Email        string
	PasswordHash string
	Role         core.Role
	RoleID       core.RoleID
}

// Contract that the Auth Domain expect to read from the User Domain
type UserReaderInterface interface {
	GetByEmailWithRoles(context.Context, string) (*AuthUser, error)
	ProvisionGuest(ctx context.Context) (*AuthUser, error)
}
