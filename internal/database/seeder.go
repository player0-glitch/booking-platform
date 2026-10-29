package database

import (
	user "booking-platform/internal/modules/user"
	"context"
)

/*
permissions
INSERT INTO permissions (permission_key, description)
VALUES
    ('users.read', 'Read user records'),
    ('users.manage', 'Create and manage user records'),
    ('roles.read', 'Read role definitions'),
    ('roles.manage', 'Create and update role definitions'),
    ('roles.grant', 'Grant and revoke roles'),
    ('organizations.read', 'Read organization data'),
    ('organizations.manage', 'Manage organization data'),
    ('bookings.create', 'Create bookings'),
    ('bookings.read.own', 'Read own bookings'),
    ('bookings.manage', 'Manage bookings within an authorized scope'),
    ('services.read', 'Read bookable services'),
    ('services.manage', 'Manage services and availability')
ON CONFLICT (permission_key) DO NOTHING;

// permissions to roles
-- 7. Seed initial role-permission mappings.
-- These are examples. Adjust them to your actual authorization policy.
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.role_name = 'Super Admin'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p
  ON p.permission_key IN (
      'organizations.read',
      'organizations.manage',
      'bookings.read.own',
      'bookings.create',
      'bookings.manage',
      'services.read',
      'services.manage'
  )
WHERE r.role_name = 'Admin'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p
  ON p.permission_key IN (
      'bookings.create',
      'bookings.read.own',
      'services.read'
  )
WHERE r.role_name = 'Customer'
ON CONFLICT DO NOTHING;

*/

type UserSeederInterface interface {
	SeedUserIfNotExist(context.Context, *user.User, int) error
	BatchSeedUser(context.Context, []user.User) error
}

type Seeder struct {
	userRepo UserSeederInterface
}

func NewSeeder(userRepo UserSeederInterface) *Seeder {
	return &Seeder{
		userRepo: userRepo,
	}
}
func (s *Seeder) Seed(ctx context.Context) error {
	return seedUsers(ctx, s.userRepo)
}
