package core

type Role string

const (
	RoleSuperAdmin Role = "SuperAdmin"
	RoleAdmin      Role = "Admin"
	RoleCustomer   Role = "Customer"
	RoleGuest      Role = "Guest"
)
