package user

import (
	"booking-platform/internal/core"
	"context"
	"errors"

	Logger "github.com/sirupsen/logrus"
)

type rolesService struct {
	roleRepo RolesRepository
}

func NewRolesService(r RolesRepository) *rolesService {
	return &rolesService{
		roleRepo: r,
	}
}

// godoc
// This assigns a role to the app's user. This function is destructive and
// and overwrites the role the user already has with the new role
func (s *rolesService) ChangeUserRole(ctx context.Context, dto *userRequest, roleName core.Role) error {

	//clear our roles and assign keys to the correct one
	// s.clearRoles(dto)
	newRoles := Role{
		RoleName: string(roleName),
		// Model:    gorm.Model{ID: uint(roleId)},
	}
	Logger.Debugf("Newly Created Role: %+v", newRoles)
	user := &User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Role:      newRoles,
	}
	return s.roleRepo.ChangeUserRole(ctx, user, &newRoles)
}

// godoc this is unstable. Delete role from a user and set them to the
// lowest level, Guest
func (s *rolesService) RemoveRoleFromUser(ctx context.Context,
	user *userRequest, role core.Role) error {
	return errors.New("Not Implemented")
}
func (s *rolesService) clearRoles(u *userRequest) {
	if u == nil {
		return
	}
	//zero out all the elements
	// clear(u.Roles)
	//reslice length back to 0 for efficiency
	// u.Roles = u.Roles[:0]
}
