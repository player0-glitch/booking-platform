package sessions

import (
	"booking-platform/internal/core"
	"time"

	"github.com/google/uuid"
)

// avoiding naming conflicts
type AppSession struct {
	core.Base
	//I might end up having to map the sessions.ID to the tokenHash
	TokenHash       []byte
	Data            []byte     `gorm:"type:jsonb;not null"`
	ClaimedByUserId *uuid.UUID `gorm:"type:uuid"`
	ExpiresAt       time.Time
}

func (AppSession) TableName() string {
	return "app_sessions"
}
