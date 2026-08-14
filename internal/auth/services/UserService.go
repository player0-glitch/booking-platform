package services

import (
	"booking-platform/internal/auth/repositories"
)

type CreateUserInput struct {
	Name     string
	LastName string
	Password string
}

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepository,
	}
}
