package user

import (
	"context"
)

type RolesRepoInterface interface {
	RemoveRoleFromUser(ctx context.Context, user *User) error
	ChangeUserRole(ctx context.Context, user *User) error
}
type rolesService struct {
	roleRepo RolesRepoInterface
}

func NewRolesService(r RolesRepoInterface) *rolesService {
	return &rolesService{
		roleRepo: r,
	}
}
