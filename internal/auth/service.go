package auth

import (
	"booking-platform/internal/database"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	errInvalidCredentials = errors.New("JWT Credentials Are Invalid")
	errAccountLocked      = errors.New("Account locked for 60 seconds due to multiple failed login attempts.")
	errInvalidToken       = errors.New("JWT Token Are Invalid")
)

type passwordRepoInterface interface {
	FindUserToken(ctx context.Context, email string) (uuid.UUID, error)
	LookupToken(ctx context.Context, tokenHash string) (*PasswordReset, error)
	DeleteEntry(ctx context.Context, tokenHash string) error
	CreateEntry(ctx context.Context, userId uuid.UUID, hash string, expiresAt time.Time) error
	UpdatePasswordReset(ctx context.Context, hashPassword []byte, userID uuid.UUID) error
}

type AuthService struct {
	userReader   UserReaderInterface //Contract that reads from the user
	jwt          jwtTokenValidator
	store        *storeManager
	passwordRepo passwordRepoInterface
	//we search the user by their email
	loginInfo  *concurrentLoginAttemptStore
	transactor database.TransactorInterface
}

func NewAuthService(
	UserReaderInterface UserReaderInterface, jwtInterface jwtTokenValidator,
	storeManager *storeManager, pr passwordRepoInterface,
	trx database.TransactorInterface) *AuthService {

	info := NewConcurrentLoginAttemptStore()

	return &AuthService{
		userReader:   UserReaderInterface,
		jwt:          jwtInterface,
		store:        storeManager,
		loginInfo:    info,
		passwordRepo: pr,
		transactor:   trx,
	}
}

// godoc
// Logs in a user allowing 3 login attempts before temporarily locking the
// user out
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	if email == "" || password == "" {
		return "", fmt.Errorf("Provide an Email and Password")
	}

	user, err := s.userReader.GetByEmail(ctx, email)
	// could not get the user based on email
	// put the error message in the string returned
	if err != nil {
		return "Invalid Email, No Database Match", err
	}

	info := s.loginInfo.GetOrAdd(user.Email)
	fmt.Printf("Login Info Dictionanry: %v\n", info)
	//Email currently locked out
	if info.lockoutEnd.After(time.Now()) {
		return "Email currently locked out", nil
	}

	//Lockout timed out but user wasn't asked to reset password yet
	if info.lockoutEnd.Before(time.Now()) &&
		info.attempts >= 3 && !info.requiresReset {
		info.requiresReset = true
		return "", fmt.Errorf("Account Locked, Reset Your Account")
	}
	//failed to login
	if email != user.Email ||
		!checkPassword(password, user.PasswordHash) {
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
	return s.jwt.CreateNewToken(user.Id)
}

// TODO Eventually remove string return
func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	//GUESTs cannot request to reset passwords
	userID, err := s.passwordRepo.FindUserToken(ctx, email)
	if err != nil {
		return "", fmt.Errorf("We Have Not Send You A Password Reset: %v", err)
	}
	//Create the token
	randomBytes := make([]byte, 32)
	_, err = rand.Read(randomBytes)

	if err != nil {
		return "", fmt.Errorf("Failed Creating Bytes For Tokens: %v", err)
	}

	rawToken := hex.EncodeToString(randomBytes)
	hashBytes := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hashBytes[:])
	expiresAt := time.Now().Add(time.Minute * 10)

	err = s.passwordRepo.CreateEntry(ctx, userID, tokenHash, expiresAt)

	if err != nil {
		return "", err
	}
	info := s.loginInfo.GetOrAdd(email)
	info.requiresReset = true
	// TODO: attatch an email service in here
	fmt.Printf("[ Email Service ] Reset Token Sent to [%s] With Token = %s ",
		email, rawToken)

	return rawToken, nil
}

// godoc
// Function used to reset the password that the reset request was sent for
func (s *AuthService) ResetPassword(ctx context.Context, email, password, rawToken string) error {
	//Only do this if a user requested for a reset
	info := s.loginInfo.GetOrAdd(email)
	if !info.requiresReset {
		return fmt.Errorf("This user did not request a password reset")
	}

	//hash the plaintoken
	hashByte := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hashByte[:])

	record, err := s.passwordRepo.LookupToken(ctx, tokenHash)
	if err != nil {
		return errors.New("Invalid Token Or Expired Token. Request A New Reset")
	}

	//delete expired tokens
	if time.Now().After(record.ExpiresAt) {
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
	err = s.transactor.Transaction(ctx, func(txCtx context.Context) error {
		//Remove this user from the block list
		s.loginInfo.Remove(email)
		//update the user's password hash
		err = s.passwordRepo.UpdatePasswordReset(txCtx, hashPassword, record.UserID)
		if err != nil {
			//triggers a rollback
			return err
		}
		//delete the token from the db
		err = s.passwordRepo.DeleteEntry(txCtx, record.TokenHash)
		if err != nil {
			//triggers a rollback
			return err
		}

		return nil //automatically commit when done
	})
	return nil
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}
