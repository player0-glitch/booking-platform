package authorisation

import (
	"context"
	"errors"
	"time"

	"booking-platform/internal/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) RevokeRole(ctx context.Context, grantID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&RoleGrant{}).
		Where("id = ? AND revoked_at IS NULL", grantID).
		Update("revoked_at", time.Now().UTC()).Error
}

func (r *Repository) HasPermissions(
	ctx context.Context,
	userID uuid.UUID,
	permissionKey string,
	organisationID *uuid.UUID,
) (bool, error) {
	// I think for performance reasons and since we're not returning
	// any data.Select (1) should suffice our exisistence check
	var count int64
	dbQuery := r.db.WithContext(ctx).
		Table("user_role_grants AS urg").
		Select("1").
		Joins("JOIN roles AS r ON r.id = urg.role_id").
		Joins("JOIN role_permissions AS rp ON rp.role_id = r.id").
		Joins("JOIN permissions AS p ON p.id = rp.permission_id ").
		Joins("JOIN users AS u ON u.id = urg.user_id").
		Where("urg.user_id = ?", userID).
		Where("p.permission_key = ?", permissionKey).
		Where("urg.revoked_at IS NULL").
		Where("(urg.expires_at IS NULL OR urg.expires_at > ?", time.Now()).
		Where("r.delete_at IS NULL").
		Where("u.delete_at IS NULL").
		Where("u.account_status = ? ", "registered")

	if organisationID == nil {
		dbQuery = dbQuery.Where("urg.organisation_id IS NULL")
	} else {
		dbQuery = dbQuery.Where("(urg.organisation_id IS NULL OR urg.organisation_id)",
			*organisationID)
	}

	// Executing the count (instead of checking for all of the possible returns,
	// we're only checking for atleast 1)
	if err := dbQuery.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) CreateGrant(ctx context.Context, grant *RoleGrant) error {
	db := database.DBFromTransaction(ctx, r.db)
	return db.Create(grant).Error
}

func (r *Repository) CheckUserRegistration(ctx context.Context,
	userID uuid.UUID,
	status string,
) error {
	user := struct {
		ID            uuid.UUID
		AccountStatus string
	}{
		ID:            userID,
		AccountStatus: status,
	}
	db := database.DBFromTransaction(ctx, r.db)
	err := db.Table("users").Select("id, account_status").
		Where("id = ? AND delete_at IS NULL", userID).
		Take(&user).Error
	if err != nil {
		return err
	}
	if user.AccountStatus != "registered" {
		return ErrInvalidGrant
	}
	return nil
}

func (r *Repository) GetRoleByID(ctx context.Context, roleID int32) (*Role, error) {
	var role Role
	db := database.DBFromTransaction(ctx, r.db)
	result := db.Where("id = ? AND delete_at IS NULL", roleID).
		Take(&role)
	return &role, result.Error
}

func (r *Repository) PreventDuplicateActiveGrants(
	ctx context.Context, userID uuid.UUID, roleID int32, organisationID *uuid.UUID,
) error {
	var existing int64
	db := database.DBFromTransaction(ctx, r.db)
	query := db.Model(&RoleGrant{}).Where("user_id = ? AND role_id = ?",
		userID, roleID).Where("revoked_at IS NULL").
		Where("(expires_at IS NULL OR expires_at > ? )", time.Now())
	if organisationID == nil {
		query = query.Where("organisation_id IS NULL")
	} else {
		query = query.Where("organisation_id = ?", *organisationID)
	}

	if err := query.Count(&existing).Error; err != nil {
		return err
	}
	if existing > 0 {
		return ErrInvalidGrant
	}
	return nil
}

func (r *Repository) GetRoleToRevoke(
	ctx context.Context, actorID, grantID uuid.UUID,
) (*RoleGrant, error) {
	var grant RoleGrant
	if err := r.db.WithContext(ctx).Where("id = ? AND revoked_at IS NULL", grantID).
		Take(&grant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGrantNotFound
		}
		return nil, err
	}
	return &grant, nil
}
