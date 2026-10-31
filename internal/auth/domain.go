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
type AppSession struct {
	core.Base
	//I might end up having to map the sessions.ID to the tokenHash
	// TokenHash       []byte
	Data            []byte     `gorm:"type:jsonb;not null"`
	ClaimedByUserId *uuid.UUID `gorm:"type:uuid"`
	ExpiresAt       time.Time
}

func (AppSession) TableName() string {
	return "app_sessions"
}
