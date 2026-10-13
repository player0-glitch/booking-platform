package database

import (
	"booking-platform/internal/core"
	"booking-platform/internal/modules/user"
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	superAdminPassword = os.Getenv("PASSWORD_SUPER_ADMIN")
	adminPassword      = os.Getenv("PASSWORD_ADMIN")
	customerPassword   = os.Getenv("PASSWORD_CUSTOMER")
	guestPassword      = ""
)

type roles struct {
	Id       int
	RoleName core.Role
}
type seedUser struct {
	FirstName string
	LastName  string
	Email     string
	Role      []roles
	// RoleID      int
	PasswordEnv string
}

var userSeeds = []seedUser{
	{"Super", "User", "super@example.com", []roles{{Id: 1, RoleName: core.RoleSuperAdmin}}, superAdminPassword},
	{"Admin", "User", "admin@example.com", []roles{{Id: 2, RoleName: core.RoleAdmin}}, adminPassword},
	{"Customer", "User", "customer@example.com", []roles{{Id: 3, RoleName: core.RoleCustomer}}, customerPassword},
	{"Guest", "User", "guest@example.com", []roles{{Id: 4, RoleName: core.RoleGuest}}, guestPassword},
}

func SeedUsers(ctx context.Context, userRepo user.UserRepository) error {
	passwords := make(map[string]string, len(userSeeds))
	var missing []string
	//Ensure existence
	for i := 0; i < len(userSeeds)-1; i++ {
		pw := userSeeds[i]
		if strings.TrimSpace(pw.PasswordEnv) == "" {
			missing = append(missing, pw.PasswordEnv)
			continue
		}
		passwords[userSeeds[i].PasswordEnv] = pw.PasswordEnv
	}
	if len(missing) > 0 {
		return fmt.Errorf("Missing Seed Passwords in .env: %s", strings.Join(missing, ", "))
	}

	//hash the passwords (all besides the Guest)
	hashes := make(map[string]string, len(userSeeds))
	for _, s := range userSeeds {
		b, err := bcrypt.GenerateFromPassword([]byte(passwords[s.PasswordEnv]), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("Failed hashing password for %s: %w",
				s.Email, err)
		}
		hashes[s.PasswordEnv] = string(b)
	}

	// insert into db
	for _, s := range userSeeds {
		hash := hashes[s.PasswordEnv]
		u := &user.User{
			FirstName:    s.FirstName,
			LastName:     s.LastName,
			Email:        s.Email,
			PasswordHash: hash,
		}
		for _, r := range s.Role {
			ur := &user.Role{
				Id:       uint(r.Id),
				RoleName: string(r.RoleName),
			}
			u.Roles = append(u.Roles, *ur)
		}
	}
	return nil

}

// I think this is a generic
// func prt[T any](v T) *T { return &v }
