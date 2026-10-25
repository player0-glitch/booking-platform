package auth

import (
	"booking-platform/internal/core"
	"booking-platform/internal/core/authorisation"
	"net/http"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type authenticatedUser struct {
	id   uuid.UUID
	role core.Role //defined as a type in auth service file
}
type Claims struct {
	UserId uuid.UUID   `json:"user_id"`
	Role   core.Role   `json:"role"`
	RoleID core.RoleID `json:"role_id"`
	jwt.RegisteredClaims
}

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
		//save the token in here
		if err != nil || tokenStr == "" || claims == nil {
			claims, err = m.service.CreateProvisionGuest(r, w)
			if err != nil {
				http.Error(w, "Cannot Create A Provisioned Guest", http.StatusInternalServerError)
				return
			}
		}

		// claims, err := m.service.ValidateToken(tokenStr)
		// if err != nil {
		// 	http.Error(w, "Unauthorised: No Valid Token Or It Expired", http.StatusUnauthorized)
		// 	return
		// }

		// authUser := authenticatedUser{
		// 	id:   claims.UserId,
		// 	role: claims.Role,
		// }

		principal := authorisation.Principal{
			UserID: claims.UserId.String(),
			Role:   claims.Role,
		}
		// ctx := context.WithValue(r.Context(), authorisation.UserContextKey, authUser)
		ctx := authorisation.WithPrincipal(r.Context(), principal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//godoc
//Ensures every visitor (especially anonymouse) has a temp database record and a valid session token
