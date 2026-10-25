package auth

import (
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
