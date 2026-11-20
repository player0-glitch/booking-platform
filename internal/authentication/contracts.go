package authentication

import (
	"context"

	"github.com/google/uuid"
)

type AuthUser struct {
	ID            uuid.UUID
	Email         string
	PasswordHash  *string
	AccountStatus string
}

type UserReaderInterface interface {
	GetByEmail(ctx context.Context, email string) (*AuthUser, error)
}
