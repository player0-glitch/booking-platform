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
