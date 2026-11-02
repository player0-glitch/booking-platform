package seeder

import (
	"booking-platform/internal/database"
	user "booking-platform/internal/modules/user"
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrFailedUserSeed            = errors.New("Failed To Seed Users")
	ErrFailedRoleSeed            = errors.New("Failed To Seed Roles")
	ErrFailedPermissionSeed      = errors.New("Failed To Seed Permissions")
	ErrFailedRolePermissionsSeed = errors.New("Failed To Seed Role Permissions Table")
)

type UserSeederInterface interface {
	BatchSeedUser(context.Context, []user.User) error
}

type SeederParams struct {
	UserRepo    UserSeederInterface
	Transaction database.TransactorInterface
}
type Seeder struct {
	DB       *gorm.DB
	userRepo UserSeederInterface
	Trx      database.TransactorInterface
}

func NewSeeder(db *gorm.DB, params SeederParams) *Seeder {
	return &Seeder{
		DB:       db,
		userRepo: params.UserRepo,
		Trx:      params.Transaction,
	}
}
func (s *Seeder) Seed(ctx context.Context) error {
	if err := s.seedUsers(ctx); err != nil {
		return ErrFailedUserSeed
	}
	/*
			Roles need to be seeded first, followed by Permissions
			and lastly the joining table RolePermissions
		   (1)
	*/
	if err := s.seedRoles(ctx); err != nil {
		return ErrFailedRoleSeed
	}
	// (2)
	if err := s.seedPermission(ctx); err != nil {
		return ErrFailedPermissionSeed
	}
	// (3)
	if err := s.seedRolePermissions(ctx); err != nil {
		return ErrFailedRolePermissionsSeed
	}
	return nil
}
