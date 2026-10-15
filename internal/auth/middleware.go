// Validation Middleware lives here
package auth

import (
	"context"
	"net/http"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserId uuid.UUID `json:"user_id"`
	Role   UserRole  `json:"role"`
	jwt.RegisteredClaims
}

// contextKey defined in context.go in this package
const UserContextKey contextKey = "user_claims"

// godoc
// Authentication Layer using jwt Tokens from claims
func (m *Module) AuthMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := m.store.GetToken(r)

		var claims *Claims

		if err == nil || tokenStr != "" {
			claims, err = m.service.ValidateToken(tokenStr)
			// http.Error(w, "Unauthorised: No Active App Session", http.StatusUnauthorized)
			// return
		}
		//No token or claim exists in the session
		//Instead of Unauthorising a user, provision them a 'guest' role
		if err != nil || tokenStr == "" || claims == nil {
			claims, err := m.service.CreateProvisionGuest(r, w)
			if err != nil {
				http.Error(w, "Internal Auth Server Error", http.StatusInternalServerError)
				return
			}
		}
		// claims, err := m.service.ValidateToken(tokenStr)
		// if err != nil {
		// 	http.Error(w, "Unauthorised: No Valid Token Or It Expired", http.StatusUnauthorized)
		// 	return
		// }
		authUser := AuthenticatedUser{
			Id:   claims.UserId,
			Role: claims.Role,
		}
		ctx := context.WithValue(r.Context(), UserContextKey, authUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//godoc
//Ensures every visitor (especially anonymouse) has a temp database record and a valid session token
