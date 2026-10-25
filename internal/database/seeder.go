package database

import (
	"booking-platform/internal/core"
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

type UserSeederInterface interface {
	SeedUserIfNotExist(context.Context, *user.User, int) error
	BatchSeedUser(context.Context, []user.User) error
}

type seedUser struct {
	FirstName   string
	LastName    string
	Email       string
	RoleID      int
	PasswordEnv string
}
type Seeder struct {
	userRepo UserSeederInterface
}

func NewSeeder(userRepo UserSeederInterface) *Seeder {
	return &Seeder{
		userRepo: userRepo,
	}
}

var userSeeds = []seedUser{
	{"Super", "User", "super@example.com" /*role{RoleName: core.RoleSuperAdmin},*/, int(core.RoleSuperAdminID), superAdminPassword},
	{"Admin", "User", "admin@example.com" /*role{RoleName: core.RoleAdmin},*/, int(core.RoleAdminID), adminPassword},
	{"Customer", "User", "customer@example.com" /* role{RoleName: core.RoleCustomer},*/, int(core.RoleCustomerID), customerPassword},
	{"Guest", "User", "guest@example.com" /*role{RoleName: core.RoleGuest},*/, int(core.RoleGuestID), guestPassword},
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

	userToInsert := make([]user.User, len(userSeeds))
	for i, s := range userSeeds {
		hash := hashes[s.PasswordEnv]
		u := user.User{
			FirstName:    s.FirstName,
			LastName:     s.LastName,
			Email:        s.Email,
			PasswordHash: hash,
			RoleId:       s.RoleID,
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
