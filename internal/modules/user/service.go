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

// public beyond the user package
// User Service authinterface, implicity implemented by userService
type IUserService interface {
	CreateUser(ctx context.Context, name, lastName, email, password string) (*userRequest, error)
	GetByEmail(ctx context.Context, email string) (*userRequest, error)
	GetById(ctx context.Context, id uint) (*userRequest, error)
	SoftDelete(ctx context.Context, id uint) error
	DeleteById(ctx context.Context, id uint) error
	FindAll(ctx context.Context) ([]userRequest, error)
}

// Concrete implementation of IUserService private to the user package
type userService struct {
	userRepo  UserRepository
	rolesRepo RolesRepository
}

func NewUserService(userRepository UserRepository, rolesRepository RolesRepository) *userService {
	return &userService{
		userRepo:  userRepository,
		rolesRepo: rolesRepository,
	}
}

func (u *userService) CreateUser(ctx context.Context, name, lastName, email, password string) (*userRequest, error) {

	passwordHash, err := u.hashPassword(password)
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
	dto := newUserRequest(user)
	return &dto, nil
}

func (s *userService) GetByEmail(ctx context.Context,
	email string) (*userRequest, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	dto := newUserRequest(user)
	return &dto, nil
}

func (s *userService) GetById(ctx context.Context, id uint) (*userRequest, error) {
	user, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		//idk how to check for this kind of error
		return nil, err
	}
	dto := newUserRequest(user)
	return &dto, nil
}

func (s *userService) DeleteById(ctx context.Context, id uint) error {
	return s.userRepo.DeleteById(ctx, id)
}

func (s *userService) SoftDelete(ctx context.Context, id uint) error {
	return s.userRepo.SoftDeleteById(ctx, id)
}

func (s *userService) FindAll(ctx context.Context) ([]userRequest, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		fmt.Printf("A User Model Array Couldn't Be formed: %x", err)
		return nil, err
	}

	fmt.Printf("A User Model: %+v", users[0])
	reqs := make([]userRequest, 0, len(users))
	for _, u := range users {
		reqs = append(reqs, newUserRequest(&u))
	}
	return reqs, nil
}

func (s *userService) hashPassword(password string) (string, error) {
	//automatic salt generation happens under the hood
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to generate hashed password: %w", err)
	}
	return string(bytes), nil
}
