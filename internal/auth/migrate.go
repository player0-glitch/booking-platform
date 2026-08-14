// DO NOT USE THIS
package auth

import (
	"booking-platform/internal/auth/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
	)
}
