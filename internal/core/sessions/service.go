package sessions

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrInvalidSession = errors.New("invalid anonymous session")
	ErrSessionClaimed = errors.New("anonymous session already claimed")
)

type Service struct {
	db  *gorm.DB
	ttl time.Duration
}

func NewService(DB *gorm.DB, ttl time.Duration) *Service {
	return &Service{
		db:  DB,
		ttl: ttl,
	}
}

// session token used for anonymous browsing and should be saved as the app session
// for persistence
func newToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// godoc
// hash the token for the the database
func hashToken(token string) []byte {
	sum := sha512.Sum512_256([]byte(token))
	return sum[:]
}
