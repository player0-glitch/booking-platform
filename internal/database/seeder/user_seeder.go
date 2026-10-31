package seeder

import (
	user "booking-platform/internal/modules/user"
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

type seedUser struct {
	FirstName string
	LastName  string
	Email     string
	// RoleID      int
	PasswordEnv string
}

var userSeeds = []seedUser{
	{"Super", "User", "super@example.com" /*role{RoleName: core.RoleSuperAdmin},*/, superAdminPassword},
	{"Admin", "User", "admin@example.com" /*role{RoleName: core.RoleAdmin},*/, adminPassword},
	{"Customer", "User", "customer@example.com" /* role{RoleName: core.RoleCustomer},*/, customerPassword},
	{"Guest", "User", "guest@example.com" /*role{RoleName: core.RoleGuest},*/, guestPassword},
}

func seedUsers(ctx context.Context, userRepo UserSeederInterface) error {
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

	userToInsert := make([]user.User, len(userSeeds))
	for i, s := range userSeeds {
		hash := hashes[s.PasswordEnv]
		u := user.User{
			FirstName:    &s.FirstName,
			LastName:     &s.LastName,
			Email:        &s.Email,
			PasswordHash: &hash,
			// Role:         user.Role{Id: s.RoleID},
		}
		//n+1
		userToInsert[i] = u
		// userRepo.SeedUserIfNotExist(ctx, &u, u.RoleId)
	}
	if err := userRepo.BatchSeedUser(ctx, userToInsert); err != nil {
		return fmt.Errorf("Failed To Seed Users: %v", err.Error())
	}
	return nil
}
