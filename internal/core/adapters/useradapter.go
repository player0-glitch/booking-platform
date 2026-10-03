package adapters

import (
	"booking-platform/internal/auth"
	"booking-platform/internal/modules/user"
	"context"
)

type UserAdapter struct {
	userRepo *user.UserRepository
}

func NewUserAdapter(userRepo *user.UserRepository) *UserAdapter {
	return &UserAdapter{
		userRepo: userRepo,
	}
}

func (a *UserAdapter) GetByEmailWithRoles(ctx context.Context, email string) (*auth.AuthUser, error) {
	userModel, err := a.userRepo.GetByEmailWithRoles(ctx, email)
	if err != nil {
		return nil, err
	}

	authUser := &auth.AuthUser{
		Id:           userModel.ID,
		Email:        userModel.Email,
		PasswordHash: userModel.PasswordHash,
		Roles:        make([]auth.UserRole, 0, len(userModel.Roles)),
	}
	for _, role := range userModel.Roles {
		authUser.Roles = append(authUser.Roles, auth.UserRole{
			Id:   role.ID,
			Name: role.RoleName,
		})
	}
	return authUser, nil
}

// If the contract between the Auth and User module this will throw
// a build error
var _ auth.UserReader = (*UserAdapter)(nil)
