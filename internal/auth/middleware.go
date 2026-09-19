// Validation Middleware lives here
package auth

import (
	"context"
	"net/http"
)

type contextKey string

const UserContextKey contextKey = "user_claims"

func (m *Module) AuthMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := m.store.GetToken(r)

		if err != nil || tokenStr == "" {
			http.Error(w, "Unauthorised: No Active App Session", http.StatusUnauthorized)
			return
		}
		claims, err := m.service.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorised: No Valid Token Or It Expired", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
