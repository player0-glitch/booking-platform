package user

import (
	"booking-platform/internal/user/models"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserInput struct {
	Name     string
	LastName string
	Password string
}

type UserService struct {
	userRepo *UserRepository
}

var (
	ErrUserNotFound = errors.New("User Not Found")
)

func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{
		userRepo: userRepository,
	}
}

// all these parameters are types
func (u *UserService) CreateUser(ctx context.Context, name, lastName, email, passsword string) (*models.User, error) {

	passwordHash, err := u.hashPassword(passsword)
	if err != nil {
		return nil, fmt.Errorf("Password hashing failed: %w", err)
	}
	user := &models.User{
		Email:        email,
		FirstName:    name,
		LastName:     lastName,
		PasswordHash: passwordHash,
	}
	notFoundError := u.userRepo.Create(ctx, user)
	if notFoundError != nil {
		return nil, notFoundError
	}
	return user, nil
}

func (s *UserService) GetById(ctx context.Context, id int) (*models.User, error) {
	user, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) DeleteById(ctx context.Context, id int) error {
	return s.userRepo.DeleteById(ctx, id)
}

func (s *UserService) SoftDelete(ctx context.Context, id int) error {
	return s.userRepo.SoftDeleteById(ctx, id)
}
func (s *UserService) FindAll(ctx context.Context) ([]models.User, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (s *UserService) hashPassword(password string) (string, error) {
	//automatic salt generation happens under the hood
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to generate hashed password: %w", err)
	}
	return string(bytes), nil
}

func (s *UserService) checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}
