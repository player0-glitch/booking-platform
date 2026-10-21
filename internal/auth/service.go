package auth

import (
	"booking-platform/internal/core"
	"context"
	"errors"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("JWT Credentials Are Invalid")
	ErrInvalidToken       = errors.New("JWT Token Are Invalid")
)

type AuthService struct {
	jwtSecret  []byte
	userReader UserReader //Contract that reads from the user
	store      *StoreManager
}

func NewAuthService(userReader UserReader, jwtSecret []byte, storeManager *StoreManager) *AuthService {

	return &AuthService{
		jwtSecret:  jwtSecret,
		userReader: userReader,
		store:      storeManager,
	}
}

func (s *AuthService) Authenticate(ctx context.Context, email, password string) (string, error) {
	user, err := s.userReader.GetByEmailWithRoles(ctx, email)
	// could not get the user based on email
	// put the error message in the string returned
	if err != nil {
		return "Invalid Email, No Database Match", err
	}

	if email != user.Email ||
		!s.checkPassword(password, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}
	//Query DB to get the role from the user.roleId to see what role this
	//user has
	claims := Claims{
		UserId: user.Id,
		// Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role: user.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString,
		&Claims{}, func(token *jwt.Token) (interface{}, error) {
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

func (s *AuthService) checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func (s *AuthService) CreateProvisionGuest(r *http.Request, w http.ResponseWriter) (*Claims, error) {
	user, err := s.userReader.ProvisionGuest(r.Context())
	if err != nil {
		return nil, err
	}
	claims := &Claims{
		UserId: user.Id,
		Role: UserRole{
			Id:   core.RoleGuestID,
			Name: string(core.RoleGuest),
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt://Guests get max 24 hour access
			jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	//create guest token and save to store
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//very unlikely to fail
	tokenStr, _ := token.SignedString(s.jwtSecret)
	s.store.SaveToken(w, r, tokenStr)

	return claims, nil
}
