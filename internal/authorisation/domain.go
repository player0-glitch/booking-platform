package authorisation

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          int32 `gorm:"primaryKey"`
	RoleName    string
	RoleKey     string `gorm:"size:50"`
	Description string
	IsSystem    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

func (Role) TableName() string {
	return "roles"
}

type Permission struct {
	ID            int32  `gorm:"primaryKey"`
	PermissionKey string `gorm:"size:120;not null"`
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}

func (Permission) TableName() string {
	return "Permissions"
}

type RolePermissions struct {
	RoleID        int32 `gorm:"primaryKey"`
	PermissionsID int32 `gorm:"primaryKey"`
	CreatedAt     time.Time
}

func (RolePermissions) TableName() string {
	return "role_permissions"
}

type RoleGrant struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID `gorm:"type:uuid;foreignKey"`
	OrganisationID uuid.UUID `gorm:"type:uuid;foreignKey"`
	RoleID         int32     `gorm:"foreignKey"`
	GrantedBy      uuid.UUID `gorm:"type:uuid;foreignKey"`
	ExpiresAt      *time.Time
	RevokedAt      *time.Time
	CreatedAt      time.Time
}

func (RoleGrant) TableName() string {
	return "user_role_grants"
}
