package auth

import (
	"booking-platform/internal/database"
	"context"
	"time"

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
	var userId uuid.UUID
	err := r.db.WithContext(ctx).Table("users").
		Where("email = ?", &email).
		Scan(&userId)
	return userId, err.Error
}

func (r *PasswordRepository) CreateEntry(ctx context.Context, userId uuid.UUID, hash string, expiresAt time.Time) error {
	p := &PasswordReset{
		UserID:    userId,
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

// godoc
// Updates the user's password hash with the hash from reset password and removes user from password reset entry
func (r *PasswordRepository) UpdatePasswordReset(ctx context.Context, passwordHash []byte, userId uuid.UUID) error {
	//Will pass r.db to run in a transaction if called within such a context
	// db := database.DbFromTransaction(ctx, r.db)
	return r.db.Table("users").
		Where("id = ?", userId).Update("password_hash", passwordHash).Error
}

// godoc
// Deletes a user's password request once password is successfully reset
func (r *PasswordRepository) DeleteEntryByUserId(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&PasswordReset{}).Error
}

// godoc
// Delete Password Reset entry from db
func (r *PasswordRepository) DeleteEntry(ctx context.Context, tokenHash string) error {

	db := database.DbFromTransaction(ctx, r.db)
	return db.Where("token = ?", tokenHash).Delete(&PasswordReset{}).Error
}
