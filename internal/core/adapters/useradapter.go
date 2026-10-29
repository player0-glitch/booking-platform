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

func (a *UserAdapter) GetByEmail(ctx context.Context, email string) (*auth.AuthUser, error) {
	userModel, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	authUser := &auth.AuthUser{
		Id:    (userModel.Base.ID),
		Email: *userModel.Email,
	}

	return authUser, nil
}

// If the contract between the Auth and User module this will throw
// a build error
var _ auth.UserReaderInterface = (*UserAdapter)(nil)
