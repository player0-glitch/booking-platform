package user

import (
	"booking-platform/internal/core"
	"context"

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
	return r.dbContext.Create(user).Error
}

func (r *UserRepository) SeedUserIfNotExist(ctx context.Context, user *User, roleId int) error {
	// err := r.dbContext.WithContext(ctx).Error
	// if err != nil {
	// 	return fmt.Errorf("%w: %w", gorm.ErrRecordNotFound, err)
	// }
	result := r.dbContext.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(user)
	//on a conflict nothing is injected into the result.RowsAffected
	return result.Error
}

func (r *UserRepository) GetById(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	err := r.dbContext.WithContext(ctx).
		First(&user, id).Error

	if err != nil {
		return nil, errRecordNotFound
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.dbContext.WithContext(ctx).
		Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, errRecordNotFound
	}
	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]User, error) {
	var users []User
	err := r.dbContext.WithContext(ctx).
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

func (r *UserRepository) BatchSeedUser(ctx context.Context, users []User) error {
	return r.dbContext.WithContext(ctx).Clauses(clause.
		OnConflict{DoNothing: true}).Create(users).Error
}
