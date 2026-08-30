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

// Create a user
func (repo *UserRepository) Create(user *models.User) error {
	//The error is returned only if it occurs
	return repo.dbContext.Create(&user).Error
}

func (repo *UserRepository) FindById(id int) (*models.User, error) {
	var user models.User
	err := repo.dbContext.First(&user, id).Error

	//error cheching before return valid ojeect
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) DeleteById(id int) error {
	var user models.User
	return repo.dbContext.Delete(&user, id).Error
}
func (repo *UserRepository) DeleteByEmail(email string) error {
	/*this permanently deletes the records instead of setting
	*  deleted_at column to a timestamp (Soft Delete)
	 */
	return repo.dbContext.Unscoped().Delete(&email).Error

}
