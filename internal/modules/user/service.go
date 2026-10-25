package user

import (
	"booking-platform/internal/core"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("User Not Found")
)

// Concrete implementation of IUserService private to the user package
type userService struct {
	userRepo  UserRepository
	rolesRepo RolesRepoInterface
}

func NewUserService(userRepository UserRepository, rolesRepository RolesRepoInterface) *userService {
	return &userService{
		userRepo:  userRepository,
		rolesRepo: rolesRepository,
	}
}

func (u *userService) CreateUser(ctx context.Context, name, lastName, email, password string) (*userRequest, error) {

	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("Password hashing failed: %w", err)
	}
	user := &User{
		Email:        email,
		FirstName:    name,
		LastName:     lastName,
		PasswordHash: passwordHash,
		RoleId:       int(core.RoleCustomerID),
		Role: Role{Id: int(core.RoleCustomerID),
			RoleName: string(core.RoleCustomer),
		},
	}
	notFoundError := u.userRepo.Create(ctx, user)
	if notFoundError != nil {
		return nil, notFoundError
	}
	dto := newUserRequest(user)
	return &dto, nil
}

func (u *userService) CreateUserWithRole(ctx context.Context, name, lastName, email, password string, roleId int) (*userRequest, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("Password hashing failed: %w", err)
	}
	user := &User{
		Email:        email,
		FirstName:    name,
		LastName:     lastName,
		PasswordHash: passwordHash,
	}
	err = u.userRepo.CreateWithRole(ctx, user, roleId)
	if err != nil {
		return nil, err
	}
	dto := newUserRequest(user)
	return &dto, nil
}
func (u *userService) GetByEmail(ctx context.Context,
	email string) (*userRequest, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	dto := newUserRequest(user)
	return &dto, nil
}

func (u *userService) GetById(ctx context.Context, id uuid.UUID) (*userRequest, error) {
	user, err := u.userRepo.GetById(ctx, id)
	if err != nil {
		//idk how to check for this kind of error
		return nil, err
	}
	dto := newUserRequest(user)
	return &dto, nil
}

func (u *userService) DeleteById(ctx context.Context, id uuid.UUID) error {
	return u.userRepo.DeleteById(ctx, id)
}

func (s *userService) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.SoftDeleteById(ctx, id)
}

func (s *userService) GetUsersByRole(ctx context.Context, roleName string) ([]userRequest, error) {

	r, err := core.ParseRole(roleName)
	if err != nil {
		return nil, fmt.Errorf("Invalid Role Name Provided")
	}

	id, _ := r.ID()

	users, err := s.userRepo.GetUsersByRoleId(ctx, id)

	if err != nil {

		return nil, fmt.Errorf("A User Model Array Couldn't Be formed: %x", err)
	}

	fmt.Printf("A User Model: %+v", users[0])

	reqs := make([]userRequest, 0, len(users))
	for _, u := range users {
		reqs = append(reqs, newUserRequest(&u))
	}
	return reqs, nil
}

func (s *userService) GetAll(ctx context.Context) ([]userRequest, error) {
	users, err := s.userRepo.GetAll(ctx)

	if err != nil {
		return nil, fmt.Errorf("A User Model Array Couldn't Be formed: %v", err)
	}

	fmt.Printf("A User Model: %+v", users[0])
	reqs := make([]userRequest, 0, len(users))
	for _, u := range users {
		reqs = append(reqs, newUserRequest(&u))
	}
	return reqs, nil
}

func hashPassword(password string) (string, error) {
	//automatic salt generation happens under the hood
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to generate hashed password: %w", err)
	}
	return string(bytes), nil
}
