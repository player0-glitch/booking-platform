package auth

import (
	"booking-platform/internal/core"
	"booking-platform/internal/database"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PasswordRepository struct {
	dbContext *gorm.DB
}

func NewPasswordRepository(db *gorm.DB) *PasswordRepository {
	return &PasswordRepository{
		dbContext: db,
	}
}

func (r *PasswordRepository) FindUserToken(ctx context.Context, email string) (uuid.UUID, error) {
	var userId uuid.UUID
	err := r.dbContext.WithContext(ctx).Table("users").
		Where("email = ? AND role != ?", &email, core.RoleGuest).
		Scan(&userId)
	return userId, err.Error
}

func (r *PasswordRepository) CreateEntry(ctx context.Context, userId uuid.UUID, hash string, expiresAt time.Time) error {
	p := &PasswordReset{
		UserID:    userId,
		TokenHash: hash,
		ExpiresAt: expiresAt,
	}
	return r.dbContext.WithContext(ctx).
		Create(&p).Error
}

func (r *PasswordRepository) LookupToken(ctx context.Context, tokenHash string) (*PasswordReset, error) {
	var p PasswordReset

	err := r.dbContext.WithContext(ctx).Model(&p).
		Where("token_hash = ?", tokenHash).
		Scan(p)
	return &p, err.Error
}

// godoc
// Updates the user's password hash with the hash from reset password and removes user from password reset entry
func (r *PasswordRepository) UpdatePasswordReset(ctx context.Context, passwordHash []byte, userId uuid.UUID) error {
	//Will pass r.dbContext to run in a transaction if called within such a context
	// db := database.DbFromTransaction(ctx, r.dbContext)
	return r.dbContext.Table("users").
		Where("id = ?", userId).Update("password_hash", passwordHash).Error
}

// godoc
// Deletes a user's password request once password is successfully reset
func (r *PasswordRepository) DeleteEntryByUserId(ctx context.Context, userID uuid.UUID) error {
	return r.dbContext.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&PasswordReset{}).Error
}

// godoc
// Delete Password Reset entry from db
func (r *PasswordRepository) DeleteEntry(ctx context.Context, tokenHash string) error {

	db := database.DbFromTransaction(ctx, r.dbContext)
	return db.Where("token = ?", tokenHash).Delete(&PasswordReset{}).Error
}
