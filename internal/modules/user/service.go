package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("User Not Found")
)

type UserService struct {
	userRepo *UserRepository
}

func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{
		userRepo: userRepository,
	}
}

func (u *UserService) CreateUser(ctx context.Context, name, lastName, email, passsword string) (*userRequest, error) {

	passwordHash, err := u.hashPassword(passsword)
	if err != nil {
		return nil, fmt.Errorf("Password hashing failed: %w", err)
	}
	user := &User{
		Email:        email,
		FirstName:    name,
		LastName:     lastName,
		PasswordHash: passwordHash,
	}
	notFoundError := u.userRepo.Create(ctx, user)
	if notFoundError != nil {
		return nil, notFoundError
	}
	return newUserRequest(user), nil
}

func (s *UserService) GetByEmail(ctx context.Context,
	email string) (*User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetById(ctx context.Context, id uint) (*userRequest, error) {
	user, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return newUserRequest(user), nil
}

func (s *UserService) DeleteById(ctx context.Context, id uint) error {
	return s.userRepo.DeleteById(ctx, id)
}

func (s *UserService) SoftDelete(ctx context.Context, id uint) error {
	return s.userRepo.SoftDeleteById(ctx, id)
}

func (s *UserService) FindAll(ctx context.Context) ([]userRequest, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return []userRequest{}, err
	}
	reqs := make([]userRequest, 0, len(users))
	for i, u := range users {
		reqs[i] = *newUserRequest(&u)
	}
	return reqs, nil
}

func (s *UserService) hashPassword(password string) (string, error) {
	//automatic salt generation happens under the hood
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to generate hashed password: %w", err)
	}
	return string(bytes), nil
}
