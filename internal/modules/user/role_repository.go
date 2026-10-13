package user

import (
	"context"

	"gorm.io/gorm"
)

type rolesRepo struct {
	dbContext *gorm.DB
}

func newRolesRepo(dbContext *gorm.DB) *rolesRepo {
	return &rolesRepo{
		dbContext: dbContext,
	}
}
func (r *rolesRepo) ChangeUserRole(ctx context.Context, user *User, role *Role) error {
	roles := make([]Role, 0, 1)
	//for now each user only has 1 role
	roles[0] = Role{RoleName: role.RoleName, Id: role.Id}
	return r.dbContext.Model(user).WithContext(ctx).Association("Roles").Replace(roles)
}

// godoc
// Remove the joining relationship between a user and a role from our system
func (r *rolesRepo) RemoveRoleFromUser(ctx context.Context, user *User, role *Role) error {
	return r.dbContext.Model(user).WithContext(ctx).Association("Roles").Delete(role)
}
