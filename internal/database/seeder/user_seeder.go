package seeder

import (
	"context"
	"fmt"
	"os"
	"strings"

	user "booking-platform/internal/modules/user"

	"golang.org/x/crypto/bcrypt"
)

var (
	superAdminPassword = os.Getenv("PASSWORD_SUPER_ADMIN")
	adminPassword      = os.Getenv("PASSWORD_ADMIN")
	customerPassword   = os.Getenv("PASSWORD_CUSTOMER")
)

type seedUser struct {
	FirstName   string
	LastName    string
	Email       string
	PasswordEnv string
}

var userSeeds = []seedUser{
	{
		FirstName:   "Super",
		LastName:    "User",
		Email:       "super@example.com",
		PasswordEnv: superAdminPassword,
	},
	{
		FirstName:   "Admin",
		LastName:    "User",
		Email:       "admin@example.com",
		PasswordEnv: adminPassword,
	},
	{
		FirstName:   "Customer",
		LastName:    "User",
		Email:       "customer@example.com",
		PasswordEnv: customerPassword,
	},
	{
		FirstName:   "Hair 1",
		LastName:    "Stylist 1",
		Email:       "stylist@nonhle.com",
		PasswordEnv: "hairStrongPass123",
	},
}

func (s *Seeder) seedUsers(ctx context.Context) error {
	passwords := make(map[string]string, len(userSeeds))
	var missing []string
	// Ensure existence
	for i := 0; i < len(userSeeds)-1; i++ {
		pw := userSeeds[i]
		if strings.TrimSpace(pw.PasswordEnv) == "" {
			missing = append(missing, pw.PasswordEnv)
			continue
		}
		passwords[userSeeds[i].PasswordEnv] = pw.PasswordEnv
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing seed passwords in .env: %s", strings.Join(missing, ", "))
	}

	// hash the passwords (all besides the Guest)
	hashes := make(map[string]string, len(userSeeds))
	for _, s := range userSeeds {
		b, err := bcrypt.GenerateFromPassword([]byte(passwords[s.PasswordEnv]), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed hashing password for %s: %w",
				s.Email, err)
		}
		hashes[s.PasswordEnv] = string(b)
	}

	usersToInsert := make([]user.User, len(userSeeds))
	for i, s := range userSeeds {
		hash := hashes[s.PasswordEnv]
		u := user.User{
			FirstName:    &s.FirstName,
			LastName:     &s.LastName,
			Email:        &s.Email,
			PasswordHash: &hash,
		}
		usersToInsert[i] = u
		//n+1
		// userRepo.SeedUserIfNotExist(ctx, &u, u.RoleId)
	}
	if err := s.userRepo.BatchSeedUser(ctx, usersToInsert); err != nil {
		return fmt.Errorf("failed to seed users: %v", err.Error())
	}
	return nil
}
