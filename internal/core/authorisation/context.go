package authorisation

import (
	"booking-platform/internal/core"
	"context"
	"net/http"
)

type contextKey string

const UserContextKey contextKey = "user_claims"

type Principal struct {
	UserID string
	Role   core.Role
}

func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, UserContextKey, principal)
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(UserContextKey).(Principal)
	return p, ok
}

// godoc
// shared method to enfore Gates on endpoints for authorisation
func RequiredRoles(allowedRoles ...core.Role) func(http.Handler) http.Handler {
	//make a map of role (a string of role names based on the Role structs passed in)
	allowed := make(map[core.Role]struct{}, len(allowedRoles))
	//this is some weird horrifc syntax
	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, ok := PrincipalFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorised User", http.StatusUnauthorized)
				return
			}

			if _, ok := allowed[principal.Role]; ok {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
