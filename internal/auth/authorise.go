package auth

import (
	"booking-platform/internal/core"
	"context"
	"net/http"
)

type AuthenticatedUser struct {
	Id    uint
	Roles []UserRole //defined as a type in auth service file
}
type contextKey string

const authenticatedUserKey contextKey = "authenticated_user"

func UserFromContext(ctx context.Context) (AuthenticatedUser, bool) {
	user, ok := ctx.Value(authenticatedUserKey).(AuthenticatedUser)
	return user, ok
}

// This is the middleware that handle authorisation
func RequiredRoles(allowedRoles ...core.Role) func(http.Handler) http.Handler {
	//make a map of role (a string of role names based on the Role structs passed in)
	allowed := make(map[core.Role]struct{}, len(allowedRoles))
	//this is some weird horrifc syntax
	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorised User", http.StatusUnauthorized)
				return
			}
			for _, userRole := range user.Roles {
				if _, ok := allowed[core.Role(userRole.Name)]; ok {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
