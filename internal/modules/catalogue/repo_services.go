package catalogue

import (
	"gorm.io/gorm"
)

type ServicesRepository struct {
	db *gorm.DB
	// trx database.TransactorInterface
}

func NewServicesRepository(db *gorm.DB) *ServicesRepository {
	return &ServicesRepository{
		db: db,
	}
}
