package user

import (
	"booking-platform/internal/user/models"
	"context"

	"gorm.io/gorm"
)

type UserRepository struct {
	dbContext *gorm.DB
}

// Treat this like constructor
func NewUserRepo(dbContext *gorm.DB) *UserRepository {
	return &UserRepository{
		dbContext: dbContext,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.dbContext.Create(&user).Error
}

func (r *UserRepository) FindById(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := r.dbContext.WithContext(ctx).First(&user, id).Error

	if err != nil {
		return nil, errRecordNotFound
	}
	return &user, nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.dbContext.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) DeleteById(ctx context.Context, id int) error {
	var user models.User
	return r.dbContext.Unscoped().WithContext(ctx).Delete(&user, id).Error
}

func (r *UserRepository) SoftDeleteById(ctx context.Context, id int) error {
	var user models.User
	return r.dbContext.WithContext(ctx).Delete(&user, id).Error
}

func (repo *UserRepository) DeleteByEmail(ctx context.Context, email string) error {
	/*this permanently deletes the records instead of setting
	*  deleted_at column to a timestamp (Soft Delete)
	 */
	return repo.dbContext.Unscoped().WithContext(ctx).Delete(&email).Error

}
