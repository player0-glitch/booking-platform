// Package jwt is a standalone package strictly for jwt with claims separated.
// jwt store (Service) should not be tampered with
package jwt

import (
	"errors"
	"time"

	claims "booking-platform/internal/claims"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid jwt token")
	ErrExpiredToken = errors.New("expired jwt token")
)

type Service struct {
	jwtSecret []byte
	issuer    string
}

func NewJWTService(jwtSecret string) *Service {
	issuer := "booking-platform-engine"
	return &Service{
		jwtSecret: []byte(jwtSecret),
		issuer:    issuer,
	}
}

func (s *Service) CreateNewToken(userID uuid.UUID) (string, error) {
	claims := claims.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        userID.String(),
			Issuer:    s.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) ValidateToken(tokenString string) (*claims.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString,
		&claims.Claims{}, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return s.jwtSecret, nil
		})

	if err != nil || !token.Valid {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*claims.Claims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
