package authorisation

import (
	"github.com/google/uuid"
	"time"
)

type Role struct {
	Id          int32 `gorm:"primaryKey"`
	RoleName    string
	RoleKey     string `gorm:"size:50"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

func (Role) TableName() string {
	return "roles"
}

type Permission struct {
	Id            int32  `gorm:"primaryKey"`
	PermissionKey string `gorm:"size:120;not null"`
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}

func (Permission) TableName() string {
	return "Permissions"
}

// I prefer to have all my tables as struct i can see in code
type RolePermissions struct {
	RoleId        int32 `gorm:"primaryKey"`
	PermissionsId int32 `gorm:"primaryKey"`
	CreatedAt     time.Time
}

func (RolePermissions) TableName() string {
	return "role_permissions"
}

type RoleGrant struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId         uuid.UUID `gorm:"type:uuid;foreignKey"`
	OrganisationId uuid.UUID `gorm:"type:uuid;foreignKey"`
	RoleId         int32     `gorm:"foreignKey"`
	GrantedBy      uuid.UUID `gorm:"type:uuid;foreignKey"`
	ExpiresAt      *time.Time
	RevokedAt      *time.Time
	CreatedAt      time.Time
}

func (RoleGrant) TableName() string {
	return "user_role_grants"
}
