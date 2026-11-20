package core

import (
	"fmt"
	"slices"
)

type Role string
type RoleID int

const (
	RoleSuperAdmin Role = "Super Admin"
	RoleAdmin      Role = "Admin"
	RoleCustomer   Role = "Customer"
	// RoleGuest      Role = "Guest"
)

const (
	RoleSuperAdminID RoleID = iota + 1 //1
	RoleAdminID                        //2
	RoleCustomerID                     //3
	RoleGuestID                        //4
)

// map of roles to their respective IDs
var RoleToID = map[Role]RoleID{
	RoleSuperAdmin: RoleSuperAdminID,
	RoleAdmin:      RoleAdminID,
	RoleCustomer:   RoleCustomerID,
	// RoleGuest:      RoleGuestID,
}

// Tied to RoleToID so these conversions never drift
var IdToRole = func() map[RoleID]Role {
	m := make(map[RoleID]Role, len(RoleToID))

	for r, id := range RoleToID {
		m[id] = r
	}
	return m
}()

// godoc
// Validate that this role exists in our map of defined Roles
// and their respective IDs
func (r Role) IsValid() bool {
	_, ok := RoleToID[r]
	return ok
}

// godoc Role -> ID
// Safely get role ID from role
func (r Role) ID() (RoleID, bool) {
	id, ok := RoleToID[r]
	return id, ok
}

// godoc ID -> Role
// Safely get role name from role ID
func (id RoleID) Role() (Role, bool) {
	r, ok := IdToRole[id]
	return r, ok
}

// godoc
// Global helper to parse string to Role type
func ParseRole(s string) (Role, error) {
	r := Role(s)
	if !r.IsValid() {
		return "", fmt.Errorf("String Not A Valid Role Name: %q", s)
	}
	return r, nil
}

// godoc
// Restricts the generic functions below to types whose underlying type
// is a string
type StringEnum interface {
	~string
}

// godoc
// Generic Validation method for enum types
func IsValidEnum[T StringEnum](value T, validOptions []T) bool {
	return slices.Contains(validOptions, value)
}
