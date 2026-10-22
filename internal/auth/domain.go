package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PasswordReset struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	TokenHash string    `gorm:"not null"`
	ExpiresAt time.Time
	CreatedAt time.Time
}

type IPasswordReset interface {
	CreateEntry(ctx context.Context, userId uuid.UUID, hash string, expiresAt time.Time) error
	LookupTokenExpiration(ctx context.Context, tokenHash string) (time.Time, error)
	DeleteEntry(ctx context.Context, tokenHash string) error
	FindUserToken(ctx context.Context, email string) (uuid.UUID, error)
}
