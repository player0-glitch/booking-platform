package seeder

import (
	"context"

	"booking-platform/internal/authorisation"
	"booking-platform/internal/database"

	"gorm.io/gorm/clause"
)

type role struct {
	id          int32
	roleKey     string
	roleName    string
	description string
	isSystem    bool
}
type permission struct {
	permissionKey string
	description   string
}

var roles = []role{
	{
		id:          1,
		roleName:    "Super Admin",
		roleKey:     "super_admin",
		description: "Platform-wide access",
		isSystem:    true,
	},
	{
		id:          2,
		roleKey:     "admin",
		roleName:    "Admin",
		description: "Admin access scoped",
		isSystem:    true,
	},
	{
		id:          3,
		roleKey:     "staff",
		roleName:    "Staff Member",
		description: "Staff Member platform access",
		isSystem:    true,
	},
	{
		id:          4,
		roleKey:     "stylist",
		roleName:    "Hair Stylist",
		description: "Belongs to 1 organisation",
		isSystem:    false,
	},
	{
		id:          5,
		roleKey:     "org_admin",
		roleName:    "organisation admin",
		description: "has dashboard  access to the organisation",
		isSystem:    false,
	},
}

var permissions = []permission{
	{"users.read", "Read user records"},
	{"users.manage", "Create and manage user records"},
	{"roles.read", "Read role definitions"},
	{"roles.manage", "Create and update role definitions"},
	{"roles.grant", "Grant and revoke roles"},
	{"organizations.read", "Read organization data"},
	{"organizations.manage", "Manage organization data"},
	{"bookings.create", "Create bookings"},
	{"bookings.read.own", "Read own bookings"},
	{"bookings.manage", "Manage bookings within an authorized scope"},
	{"services.read", "Read bookable services"},
	{"services.manage", "Manage services and availability"},
	{"organisations.inventory.manage", "Can manage their organisation's resources and inventory"},
	{"organisations.user.manage", "Can manage their organisation's users"},
}

func (s *Seeder) seedRoles(ctx context.Context) error {
	rolesToInsert := make([]authorisation.Role, len(roles))
	for k, v := range roles {
		rolesToInsert[k] = authorisation.Role{
			Id:          v.id,
			RoleName:    v.roleName,
			RoleKey:     v.roleKey,
			Description: v.description,
			IsSystem:    v.isSystem,
		}
	}
	if err := s.DB.WithContext(ctx).Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(rolesToInsert).Error; err != nil {
		return err
	}
	return nil
}

func (s *Seeder) seedPermission(ctx context.Context) error {
	permissionsToInsert := make([]authorisation.Permission, len(permissions))
	for k, v := range permissions {
		permissionsToInsert[k] = authorisation.Permission{
			PermissionKey: v.permissionKey,
			Description:   v.description,
		}
	}
	// Bulk Insert
	if err := s.DB.WithContext(ctx).Table("permissions").
		Clauses(clause.OnConflict{DoNothing: true}).Create(&permissionsToInsert).Error; err != nil {
		return err
	}
	return nil
}

func (s *Seeder) seedRolePermissions(ctx context.Context) error {
	// Querying the database directly avoid round trips and memory overhead using GORM
	// Using raw SQL is the best approach, as data does not need to be moved back and forth between
	// the database and GORM (so no need for the Object Relational Mapper). This is
	// database administrative work

	// Super Admin Gets All The Permissions
	query := `
	-- Super Admin gets all the permissions hence the Cross join
	INSERT INTO role_permissions (role_id,permission_id)
		SELECT r.id, p.id
		FROM roles r
		CROSS JOIN permissions p 
		WHERE r.role_key = 'super_admin'

		UNION ALL

	-- Admin Key
	--INSERT INTO role_permissions (role_id,permission_id)
	SELECT r.id, p.id
	FROM roles r
	JOIN permissions p ON p.permission_key IN (?,?,?,?,?,?,?)
	WHERE r.role_key = 'admin'
	
	UNION ALL

	SELECT r.id, p.id
	FROM roles r
	JOIN permissions p ON p.permission_key IN (?,?)
	WHERE r.role_key = 'org_admin'

	UNION ALL 

	-- Staff Key (could be employees in some businesses)
	--INSERT INTO role_permissions (role_id,permission_id)
	SELECT r.id, p.id
	FROM roles r
	JOIN permissions p ON p.permission_key IN (?,?,?,?,?)
	WHERE r.role_key = 'staff'

	-- Execute this only since on all Unions
	ON CONFLICT DO NOTHING;
	`
	adminKey := []any{
		"organizations.read",
		"organizations.manage",
		"bookings.read.own",
		"bookings.create",
		"bookings.manage",
		"services.read",
		"services.manage",
	}

	orgAdminKey := []any{
		"organisations.inventory.manage",
		"organisations.user.manage",
	}

	staffKey := []any{
		"bookings.create",
		"services.manage",
		"bookings.manage",
		"bookings.read.own",
		"services.read",
	}
	var args []any

	args = append(args, adminKey...)
	args = append(args, orgAdminKey...)
	args = append(args, staffKey...)

	// Bulk insert in a transaction
	err := s.Trx.Transaction(ctx, func(txCtx context.Context) error {
		db := database.DbFromTransaction(txCtx, s.DB)
		if err := db.Exec(query, args...).Error; err != nil {
			// Triggers a rollback
			return err
		}
		return nil // automatically commit when done
	})
	return err
}
