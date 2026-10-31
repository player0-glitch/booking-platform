package authorisation

import (
	"booking-platform/internal/database"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}
func (r *Repository) RevokeRole(ctx context.Context, grantId uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&RoleGrant{}).
		Where("id = ? AND revoked_at IS NULL", grantId).
		Update("revoked_at", time.Now().UTC()).Error
}

func (r *Repository) HasPermissions(
	ctx context.Context,
	userId uuid.UUID,
	permissionKey string,
	organisationId *uuid.UUID,
) (bool, error) {
	//I think for performance reasons and since we're not returning
	//any data.Select (1) should suffice our exisistence check
	var count int64 //should not be used
	dbQuery := r.db.WithContext(ctx).
		Table("user_role_grants AS urg").
		Select("1").
		Joins("JOIN roles AS r ON r.id = urg.role_id").
		Joins("JOIN role_permissions AS rp ON rp.role_id = r.id").
		Joins("JOIN permissions AS p ON p.id = rp.permission_id ").
		Joins("JOIN users AS u ON u.id = urg.user_id").
		Where("urg.user_id = ?", userId).
		Where("p.permission_key = ?", permissionKey).
		Where("urg.revoked_at IS NULL").
		Where("(urg.expires_at IS NULL OR urg.expires_at > ?", time.Now()).
		Where("r.delete_at IS NULL").
		Where("u.delete_at IS NULL").
		Where("u.account_status = ? ", "registered")

	if organisationId == nil {
		dbQuery = dbQuery.Where("urg.organisation_id IS NULL")
	} else {
		dbQuery = dbQuery.Where("(urg.organisation_id IS NULL OR urg.organisation_id)",
			*organisationId)
	}

	//Executing the count (instead of checking for all of the possible returns,
	//we're only checking for atleast 1)
	if err := dbQuery.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) CreateGrant(ctx context.Context, grant *RoleGrant) error {
	db := database.DbFromTransaction(ctx, r.db)
	return db.Create(grant).Error
}

func (r *Repository) CheckUserRegistration(ctx context.Context,
	userId uuid.UUID,
	status string) error {
	user := struct {
		Id            uuid.UUID
		AccountStatus string
	}{
		Id:            userId,
		AccountStatus: status,
	}
	db := database.DbFromTransaction(ctx, r.db)
	err := db.Table("users").Select("id, account_status").
		Where("id = ? AND delete_at IS NULL", userId).
		Take(&user).Error

	if err != nil {
		return err
	}
	if user.AccountStatus != "registered" {
		return ErrInvalidGrant
	}
	return nil
}
func (r *Repository) GetRoleById(ctx context.Context, roleId int32) (*Role, error) {
	var role Role
	db := database.DbFromTransaction(ctx, r.db)
	result := db.Where("id = ? AND delete_at IS NULL", roleId).
		Take(&role)
	return &role, result.Error
}

func (r *Repository) PreventDuplicateActiveGrants(
	ctx context.Context, userId uuid.UUID, roleId int32, organisationId *uuid.UUID,
) error {
	var existing int64
	db := database.DbFromTransaction(ctx, r.db)
	query := db.Model(&RoleGrant{}).Where("user_id = ? AND role_id = ?",
		userId, roleId).Where("revoked_at IS NULL").
		Where("(expires_at IS NULL OR expires_at > ? )", time.Now())
	if organisationId == nil {
		query = query.Where("organisation_id IS NULL")
	} else {
		query = query.Where("organisation_id = ?", *organisationId)
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
	ctx context.Context, actorId, grantId uuid.UUID,
) (*RoleGrant, error) {
	var grant RoleGrant
	if err := r.db.WithContext(ctx).Where("id = ? AND revoked_at IS NULL", grantId).
		Take(&grant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGrantNotFound
		}
		return nil, err
	}
	return &grant, nil
}
