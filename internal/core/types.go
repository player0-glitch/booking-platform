// Package core
package core

import (
	"fmt"
	"slices"
)

type (
	Role   string
	RoleID int
)

const (
	RoleSuperAdmin Role = "Super Admin"
	RoleAdmin      Role = "Admin"
	RoleCustomer   Role = "Customer"
)

const (
	RoleSuperAdminID RoleID = iota + 1 //1
	RoleAdminID                        //2
	RoleCustomerID                     //3
	RoleGuestID                        //4
)

// RoleToID map of roles to their respective IDs
var RoleToID = map[Role]RoleID{
	RoleSuperAdmin: RoleSuperAdminID,
	RoleAdmin:      RoleAdminID,
	RoleCustomer:   RoleCustomerID,
	// RoleGuest:      RoleGuestID,
}

// IDToRole Tied to RoleToID so these conversions never drift
var IDToRole = func() map[RoleID]Role {
	m := make(map[RoleID]Role, len(RoleToID))

	for r, id := range RoleToID {
		m[id] = r
	}
	return m
}()

// IsValid Validate that this role exists in our map of defined Roles
// and their respective IDs
func (r Role) IsValid() bool {
	_, ok := RoleToID[r]
	return ok
}

// ID Safely get role ID from role
// Role -> ID
func (r Role) ID() (RoleID, bool) {
	id, ok := RoleToID[r]
	return id, ok
}

// Role Safely get role name from role ID
// ID -> Role
func (id RoleID) Role() (Role, bool) {
	r, ok := IDToRole[id]
	return r, ok
}

// ParseRole Global helper to parse string to Role type
func ParseRole(s string) (Role, error) {
	r := Role(s)
	if !r.IsValid() {
		return "", fmt.Errorf("string not a valid role name: %q", s)
	}
	return r, nil
}

// StringEnum Restricts the generic functions below to types whose underlying type
// is a string
type StringEnum interface {
	~string
}

// IsValidEnum Generic Validation method for enum types
func IsValidEnum[T StringEnum](value T, validOptions []T) bool {
	return slices.Contains(validOptions, value)
}
