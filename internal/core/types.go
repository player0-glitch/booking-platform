package core

type Role string

const (
	RoleSuperAdmin Role = "SuperAdmin"
	RoleAdmin      Role = "Admin"
	RoleCustomer   Role = "Customer"
	RoleGuest      Role = "Guest"
)

var (
	RoleSuperAdminID = 1
	RoleAdminID      = 2
	RoleCustomerID   = 3
	RoleGuestID      = 4
)
