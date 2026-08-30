package services

import (
	"booking-platform/internal/auth/models"
	"booking-platform/internal/auth/repositories"
	"errors"
)

type CreateUserInput struct {
	Name     string
	LastName string
	Password string
}

type UserService struct {
	userRepo repositories.UserRepository
}

var (
	ErrUserNotFound = errors.New("User Not Found")
)

func NewUserService(userRepository repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepository,
	}
}

// all these parameters are types
func (u *UserService) CreateUser(name, lastName, email, passsword string) (*models.User, error) {

	user := &models.User{
		Email:        email,
		FirstName:    name,
		LastName:     lastName,
		PasswordHash: passsword,
	}

	notFoundError := u.userRepo.Create(user)
	if notFoundError != nil {
		return nil, notFoundError
	}
	return user, nil

}
