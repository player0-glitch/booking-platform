package auth

import (
	"booking-platform/internal/core"
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

func (PasswordReset) TableName() string {
	return "password_reset"
}

// avoiding naming conflicts
type appSession struct {
	core.Base
	// TokenHash       []byte
	Data            []byte     `gorm:"type:jsonb;not null"`
	ClaimedByUserId *uuid.UUID `gorm:"type:uuid"`
	ExpiresAt       time.Time
}

func (appSession) TableName() string {
	return "app_sessions"
}
