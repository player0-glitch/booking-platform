package auth

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtService struct {
	jwtSecret []byte
}

func NewJWTService(jwtStr string) *jwtService {
	return &jwtService{jwtSecret: []byte(jwtStr)}
}

func (s *jwtService) CreateNewToken(userId uuid.UUID) (string, error) {
	claims := claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
func (s *jwtService) ValidateToken(tokenString string) (*claims, error) {
	token, err := jwt.ParseWithClaims(tokenString,
		&claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errInvalidToken
			}
			return s.jwtSecret, nil
		})

	if err != nil || !token.Valid {
		return nil, errInvalidToken
	}

	claims, ok := token.Claims.(*claims)
	if !ok {
		return nil, errInvalidToken
	}
	return claims, nil
}
