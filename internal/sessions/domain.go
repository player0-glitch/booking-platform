package sessions

import (
	"time"

	"booking-platform/internal/core"

	"github.com/google/uuid"
)

type AppSession struct {
	base core.Base
	// I might end up having to map the sessions.ID to the tokenHash
	TokenHash       []byte
	Data            []byte     `gorm:"type:jsonb;not null"`
	ClaimedByUserID *uuid.UUID `gorm:"type:uuid"`
	ExpiresAt       time.Time
}

func (AppSession) TableName() string {
	return "app_sessions"
}
