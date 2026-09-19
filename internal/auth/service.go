package auth

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("JWT Credentials Are Invalid")
	ErrInvalidToken       = errors.New("JWT Token Are Invalid")
)

type Claims struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type AuthService struct {
	jwtSecret []byte
}

func NewAuthService(jwtSecret []byte) *AuthService {
	return &AuthService{
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Authenticate(email, password string) (string, error) {
	// figure out how to actually check the passowrd
	if email != "user@example.co.za" || password != "MyStrongPass123" {
		return "", ErrInvalidCredentials
	}

	claims := Claims{
		UserId: "user_123",
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
