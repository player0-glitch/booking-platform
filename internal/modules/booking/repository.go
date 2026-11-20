package booking

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(DB *gorm.DB) *Repository {
	return &Repository{
		db: DB,
	}
}
