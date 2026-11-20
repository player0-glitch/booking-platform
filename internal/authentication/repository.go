package authentication

import (
	"context"
	"time"

	"booking-platform/internal/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PasswordRepository struct {
	db *gorm.DB
}

func NewPasswordRepository(db *gorm.DB) *PasswordRepository {
	return &PasswordRepository{
		db: db,
	}
}

func (r *PasswordRepository) FindUserToken(ctx context.Context, email string) (uuid.UUID, error) {
	var userID uuid.UUID
	err := r.db.WithContext(ctx).Table("users").
		Where("email = ?", &email).
		Scan(&userID)
	return userID, err.Error
}

func (r *PasswordRepository) CreateEntry(ctx context.Context,
	userID uuid.UUID, hash string, expiresAt time.Time,
) error {
	p := &PasswordReset{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: expiresAt,
	}
	return r.db.WithContext(ctx).
		Create(&p).Error
}

func (r *PasswordRepository) LookupToken(ctx context.Context, tokenHash string) (*PasswordReset, error) {
	var p PasswordReset

	err := r.db.WithContext(ctx).Model(&p).
		Where("token_hash = ?", tokenHash).
		Scan(p)
	return &p, err.Error
}

// UpdatePasswordReset Updates the user's password hash with the hash from reset password and removes user from password reset entry
func (r *PasswordRepository) UpdatePasswordReset(ctx context.Context, passwordHash []byte, userID uuid.UUID) error {
	// Will pass r.db to run in a transaction if called within such a context
	db := database.DBFromTransaction(ctx, r.db)
	return db.Table("users").
		Where("id = ?", userID).Update("password_hash", passwordHash).Error
}

// DeleteEntryByUserID Deletes a user's password request once password is successfully reset
func (r *PasswordRepository) DeleteEntryByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&PasswordReset{}).Error
}

func (r *PasswordRepository) DeleteEntry(ctx context.Context, tokenHash string) error {
	db := database.DBFromTransaction(ctx, r.db)
	return db.Where("token = ?", tokenHash).Delete(&PasswordReset{}).Error
}
