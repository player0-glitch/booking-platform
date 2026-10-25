package user

import (
	"booking-platform/internal/core"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type role core.Role

// Roles Repository
type rolesRepo struct {
	dbContext *gorm.DB
}

func newRolesRepo(dbContext *gorm.DB) *rolesRepo {
	return &rolesRepo{
		dbContext: dbContext,
	}
}

type UserRepository struct {
	dbContext *gorm.DB
}

func NewUserRepo(dbContext *gorm.DB) *UserRepository {
	return &UserRepository{
		dbContext: dbContext,
	}
}

// godoc
// This creates a new user in storage defaulted to a role of Customer
func (r *UserRepository) Create(ctx context.Context, user *User) error {
	var role Role
	err := r.dbContext.WithContext(ctx).
		Where("role_name = ?", string(core.RoleCustomer)).
		First(&role).Error

	if err != nil {
		return fmt.Errorf("%w: %w ", errFailedToSaveCustomer, err)
	}
	return r.dbContext.Create(user).Error
}

// godoc
// Create a user with any role ("Super Admin", "Admin","Customer","Guest")
func (r *UserRepository) CreateWithRole(ctx context.Context, user *User, roleId int) error {
	var role Role
	err := r.dbContext.WithContext(ctx).Where("role_name = ?", &role, roleId).Error
	if err != nil {
		return fmt.Errorf("%w: %w", gorm.ErrRecordNotFound, err)
	}
	return r.dbContext.Create(user).Error
}
func (r *UserRepository) SeedUserIfNotExist(ctx context.Context, user *User, roleId int) error {
	var role Role
	err := r.dbContext.WithContext(ctx).Where("role_name = ?", &role, roleId).Error
	if err != nil {
		return fmt.Errorf("%w: %w", gorm.ErrRecordNotFound, err)
	}
	result := r.dbContext.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(user)
	//on a conflict nothing is injected into the result.RowsAffected
	return result.Error
}

func (r *UserRepository) GetById(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	err := r.dbContext.WithContext(ctx).
		Joins("Role").First(&user, id).Error

	if err != nil {
		return nil, errRecordNotFound
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.dbContext.WithContext(ctx).
		Joins("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, errRecordNotFound
	}
	return &user, nil
}

func (r *UserRepository) GetUsersByRoleId(ctx context.Context, id core.RoleID) ([]User, error) {

	var users []User
	err := r.dbContext.WithContext(ctx).
		Joins("Role").
		Where(&User{RoleId: int(id)}).
		Find(&users).Error

	if err != nil {
		return nil, err
	}
	return users, nil
}
func (r *UserRepository) GetAll(ctx context.Context) ([]User, error) {
	var users []User
	err := r.dbContext.WithContext(ctx).
		Joins("Role").
		Order("created_at DESC").
		// Limit(50).
		Find(&users).Error

	if err != nil {
		return nil, err
	}
	return users, nil
}

// Permanently delete by id but i might need to not use this
func (r *UserRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	var user User

	return r.dbContext.Unscoped().WithContext(ctx).Delete(&user, id).Error
}

func (r *UserRepository) SoftDeleteById(ctx context.Context, id uuid.UUID) error {
	var user User
	return r.dbContext.WithContext(ctx).Delete(&user, id).Error
}

// Permanently delete by email but i might need to not use this
func (r *UserRepository) DeleteByEmail(ctx context.Context, email string) error {
	/*this permanently deletes the records instead of setting
	*  deleted_at column to a timestamp (Soft Delete)
	 */
	return r.dbContext.Unscoped().WithContext(ctx).Delete(&email).Error
}

// godoc
// Used to create a temporary GUEST user
func (r *UserRepository) ProvisionGuest(ctx context.Context) (*User, error) {
	id := core.RoleToID[core.RoleGuest]
	user := &User{Role: Role{Id: int(id)}}
	err := r.dbContext.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", gorm.ErrRecordNotFound, err)
	}
	return user, err
}

func (r *UserRepository) BatchSeedUser(ctx context.Context, users []User) error {
	return r.dbContext.WithContext(ctx).Clauses(clause.
		OnConflict{DoNothing: true}).Create(users).Error
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
