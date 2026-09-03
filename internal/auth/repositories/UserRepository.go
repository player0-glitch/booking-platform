// auth/repositories/UserRepository.go
package repositories

import (
	"booking-platform/internal/auth/models"

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

func (r *UserRepository) Create(user *models.User) error {
	return r.dbContext.Create(&user).Error
}

func (r *UserRepository) FindById(id int) (*models.User, error) {
	var user models.User
	err := r.dbContext.First(&user, id).Error

	if err != nil {
		return nil, errRecordNotFound
	}
	return &user, nil
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.dbContext.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) DeleteById(id int) error {
	var user models.User
	return r.dbContext.Unscoped().Delete(&user, id).Error
}

func (r *UserRepository) SoftDeleteById(id int) error {
	var user models.User
	return r.dbContext.Delete(&user, id).Error
}

func (repo *UserRepository) DeleteByEmail(email string) error {
	/*this permanently deletes the records instead of setting
	*  deleted_at column to a timestamp (Soft Delete)
	 */
	return repo.dbContext.Unscoped().Delete(&email).Error

}
