package adapters

import (
	"booking-platform/internal/authentication"
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

func (a *UserAdapter) GetByEmail(ctx context.Context, email string) (*authentication.AuthUser, error) {
	userModel, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	authUser := &authentication.AuthUser{
		Id:            (userModel.Base.ID),
		Email:         *userModel.Email,
		PasswordHash:  userModel.PasswordHash,
		AccountStatus: string(userModel.AccountStatus),
	}

	return authUser, nil
}

// If the contract between the Auth and User module this will throw
// a build error
var _ authentication.UserReaderInterface = (*UserAdapter)(nil)
