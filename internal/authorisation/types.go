// Package authorisation only deals with authenticated users or provision a 'Guest' or Anonymous Role/session
package authorisation

type PermissionKey string

const (
	PermissionUserRead   PermissionKey = "users.read"
	PermissionUserManage PermissionKey = "users.manage"

	PermissionRoleRead   PermissionKey = "roles.read"
	PermissionRoleManage PermissionKey = "roles.manage"
	PermissionRoleGrant  PermissionKey = "roles.grant"

	PermissionOrganisationRead   PermissionKey = "organisations.read"
	PermissionOrganisationManage PermissionKey = "organizations.manage"

	PermissionBookingCreate  PermissionKey = "bookings.create"
	PermissionBookingManage  PermissionKey = "bookings.manage"
	PermissionBookingReadOwn PermissionKey = "bookings.read.own"

	PermissionServiceRead   PermissionKey = "services.read"
	PermissionServiceManage PermissionKey = "bookings.manage"

	PermissionOrganisationInventoryManage PermissionKey = "organisations.inventory.manage"
	PermissionOrganisationUserManage      PermissionKey = "organisations.user.manage"
)
