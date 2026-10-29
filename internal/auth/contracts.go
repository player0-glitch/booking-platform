package auth

import (
	"context"

	"github.com/google/uuid"
)

// type storeManagerInterface interface {
// 	SaveToken(w http.ResponseWriter, r *http.Request, token string) error
// 	Clear(w http.ResponseWriter, r *http.Request) error
// }

// User struct with the attributes expected byt the Auth Domain
// This is adapted to from the User model
type AuthUser struct {
	Id           uuid.UUID
	Email        string
	PasswordHash string
}

// Contract that the Auth Domain expect to read from the User Domain
type UserReaderInterface interface {
	GetByEmail(ctx context.Context, email string) (*AuthUser, error)
}
