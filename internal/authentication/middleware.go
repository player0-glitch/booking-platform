// package authentication this middleware does not create a user or a customer. It only creats a
// persistent session.
package authentication

import (
	"net/http"

	"booking-platform/internal/claims"
	"booking-platform/internal/jwt"
	"booking-platform/internal/sessions"
)

type JwtTokenValidator interface {
	ValidateToken(tokenString string) (*claims.Claims, error)
}

type AuthenticationMiddleware struct {
	// this is not from gorilla, it's a wrapper
	sess   *sessions.Manager
	tokens JwtTokenValidator
}

func NewAuthenticationMiddleware(session *sessions.Manager, token JwtTokenValidator) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		sess:   session,
		tokens: token,
	}
}

// godoc
// RequiresAuthentication only checks that a jwt token is required to be authenticated here
func (m *AuthenticationMiddleware) RequiresAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.sess.GetToken(r)
		if err != nil {
			http.Error(w, "unauthenticated", http.StatusUnauthorized)
			return
		}

		userClaims, err := m.tokens.ValidateToken(token)
		if err != nil {
			http.Error(w, "invalid or expired session", http.StatusUnauthorized)
			return
		}
		ctx := jwt.WithClaims(r.Context(), userClaims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
