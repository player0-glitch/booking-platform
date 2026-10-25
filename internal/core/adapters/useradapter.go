package adapters

import (
	"booking-platform/internal/auth"
	"booking-platform/internal/core"
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

func (a *UserAdapter) ProvisionGuest(ctx context.Context) (*auth.AuthUser, error) {
	userModel, err := a.userRepo.ProvisionGuest(ctx)
	if err != nil {
		return nil, err
	}

	authUser := &auth.AuthUser{
		Id:           (userModel.Base.ID),
		Email:        userModel.Email,
		PasswordHash: userModel.PasswordHash,
		Role:         core.Role(userModel.Role.RoleName),
		RoleID:       core.RoleID(userModel.RoleId),
	}

	return authUser, nil
}

func (a *UserAdapter) GetByEmailWithRoles(ctx context.Context, email string) (*auth.AuthUser, error) {
	userModel, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	authUser := &auth.AuthUser{
		Id:     (userModel.Base.ID),
		Email:  userModel.Email,
		Role:   core.Role(userModel.Role.RoleName),
		RoleID: core.RoleID(userModel.RoleId),
	}

	return authUser, nil
}

// If the contract between the Auth and User module this will throw
// a build error
var _ auth.UserReaderInterface = (*UserAdapter)(nil)
