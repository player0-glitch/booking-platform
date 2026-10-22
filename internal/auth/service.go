package auth

import (
	"booking-platform/internal/core"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"fmt"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	errInvalidCredentials = errors.New("JWT Credentials Are Invalid")
	errAccountLocked      = errors.New("Account locked for 60 seconds due to multiple failed login attempts.")
	ErrInvalidToken       = errors.New("JWT Token Are Invalid")
)

type AuthService struct {
	jwtSecret    []byte
	userReader   UserReader //Contract that reads from the user
	store        *StoreManager
	passwordRepo IPasswordReset
	//we search the user by their email
	// loginInfo map[email]loginAttempInfo
	loginInfo *concurrentLoginAttemptStore
}

func NewAuthService(userReader UserReader, jwtSecret []byte, storeManager *StoreManager, pr IPasswordReset) *AuthService {

	info := NewConcurrentLoginAttemptStore()

	return &AuthService{
		jwtSecret:    jwtSecret,
		userReader:   userReader,
		store:        storeManager,
		loginInfo:    info,
		passwordRepo: pr,
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
		return "", errInvalidCredentials
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

// godoc
// Logs in a user allowing 3 login attempts before temporarily locking the
// user out
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	if email == "" || password == "" {
		return "", fmt.Errorf("Provide an Email and Password")
	}

	user, err := s.userReader.GetByEmailWithRoles(ctx, email)
	// could not get the user based on email
	// put the error message in the string returned
	if err != nil {
		return "Invalid Email, No Database Match", err
	}

	info := s.loginInfo.GetOrAdd(user.Email)
	//Email currently locked out
	if info.lockoutEnd.Compare(time.Now()) == -1 {
		return "Email currently locked out", nil
	}

	//Lockout timed out but user wasn't asked to reset password yet
	if info.lockoutEnd.Compare(time.Now()) == 1 &&
		info.attempts >= 3 && !info.requiresReset {
		info.requiresReset = true
		return "", fmt.Errorf("Account Locked, Reset Your Account")
	}
	//failed to login
	if email != user.Email ||
		!s.checkPassword(password, user.PasswordHash) {
		info.attempts++

		if info.attempts >= 3 {
			//60 seconds to test
			info.lockoutEnd = time.Now().Add(time.Second * 60)
			info.requiresReset = false
			return "", errAccountLocked
		}
		return "", errInvalidCredentials
	}

	//successful login will remove you from the attempts store
	s.loginInfo.Remove(user.Email)

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

func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	//GUESTs cannot request to reset passwords
	userID, err := s.passwordRepo.FindUserToken(ctx, email)
	if err != nil {
		return fmt.Errorf("We Have Not Send You A Password Reset: %v", err)
	}
	//Create the token
	randomBytes := make([]byte, 32)
	_, err = rand.Read(randomBytes)

	if err != nil {
		return fmt.Errorf("Failed Creating Bytes For Tokens: %v", err)
	}

	rawToken := hex.EncodeToString(randomBytes)
	hashBytes := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hashBytes[:])
	expiresAt := time.Now().Add(time.Minute * 10)

	err = s.passwordRepo.CreateEntry(ctx, userID, tokenHash, expiresAt)

	if err != nil {
		return err
	}
	// TODO: attatch an email service in here
	fmt.Printf("[ Email Service ] Reset Token Sent to [%s] With Token = %s ",
		email, rawToken)

	return nil
}

// godoc
// Function used to reset the password that the reset request was sent for
func (s *AuthService) ResetPassword(ctx context.Context, password, rawToken string) error {
	//hash the plaintoken
	hashByte := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hashByte[:])

	expiresAt, err := s.passwordRepo.LookupTokenExpiration(ctx, tokenHash)
	if err != nil {
		return errors.New("Invalid Token Or Expired Token. Request A New Reset")
	}

	//delete expired tokens
	if time.Now().After(expiresAt) {
		err := s.passwordRepo.DeleteEntry(ctx, tokenHash)
		if err == nil {
			return errors.New("Expired Token Deleted")
		}
	}

	//hash new password the update user
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Failed To Hash Newly Requested Password")
	}

	//update should happen in a transaction. This operation has to be atomic

	return nil
}
