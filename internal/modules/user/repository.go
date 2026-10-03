package user

import (
	"context"

	"gorm.io/gorm"
)

// var _ core.UserReader = (*UserRepository)(nil)

type UserRepository struct {
	dbContext *gorm.DB
}

func NewUserRepo(dbContext *gorm.DB) *UserRepository {
	return &UserRepository{
		dbContext: dbContext,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
	return r.dbContext.Create(&user).Error
}

func (r *UserRepository) FindById(ctx context.Context, id uint) (*User, error) {
	var user User
	err := r.dbContext.WithContext(ctx).Preload("Roles").Limit(1).First(&user, id).Error

	if err != nil {
		return nil, errRecordNotFound
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.dbContext.WithContext(ctx).First(&user, email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) GetByEmailWithRoles(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.dbContext.WithContext(ctx).Preload("Roles").Where("email = ?").First(&user, email).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) FindAll(ctx context.Context) ([]User, error) {
	var users []User
	err := r.dbContext.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) DeleteById(ctx context.Context, id uint) error {
	var user User
	return r.dbContext.Unscoped().WithContext(ctx).Delete(&user, id).Error
}

func (r *UserRepository) SoftDeleteById(ctx context.Context, id uint) error {
	var user User
	return r.dbContext.WithContext(ctx).Delete(&user, id).Error
}

func (repo *UserRepository) DeleteByEmail(ctx context.Context, email string) error {
	/*this permanently deletes the records instead of setting
	*  deleted_at column to a timestamp (Soft Delete)
	 */
	return repo.dbContext.Unscoped().WithContext(ctx).Delete(&email).Error
}
