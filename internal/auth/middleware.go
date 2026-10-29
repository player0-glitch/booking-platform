package auth

import (
	"booking-platform/internal/core"
	"booking-platform/internal/core/authorisation"
	"net/http"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

type jwtTokenValidator interface {
	ValidateToken(tokenString string) (*claims, error)
	CreateNewToken(userId uuid.UUID) (string, error)
}

type authenticatedUser struct {
	id   uuid.UUID
	role core.Role //defined as a type in auth service file
}

type claims struct {
	UserId uuid.UUID   `json:"user_id"`
	Role   core.Role   `json:"role"`
	RoleID core.RoleID `json:"role_id"`
	jwt.RegisteredClaims
}
type AppMiddleware struct {
	// claims   claims
	// authUser authenticatedUser
	//from anon
	name  string
	sess  sessions.Store
	store *storeManager

	jwt jwtTokenValidator
}

func NewAppMiddleware(store *storeManager, name string, jwtInterface jwtTokenValidator) *AppMiddleware {

	return &AppMiddleware{
		sess:  store.session,
		store: store,
		name:  name,
		jwt:   jwtInterface,
	}
}

// godoc
// Authentication Layer using jwt Tokens from claims
func (am *AppMiddleware) AuthMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := am.store.GetToken(r)
		// 	if err != nil {
		// 		http.Error(w, "Cannot Create A Provisioned Guest", http.StatusInternalServerError)
		// 		return
		// 	}
		// }

		claims, err := am.jwt.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorised: No Valid Token Or It Expired", http.StatusUnauthorized)
			return
		}

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

// godoc
// Ensures every visitor (especially anonymouse) has a temp database record and a valid session token
func (am *AppMiddleware) Ensure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := am.sess.Get(r, am.name)
		if err != nil {
			http.Error(w, "Unable to load session", http.StatusInternalServerError)
			return
		}
		//Provision a guest by marking a session as anonymous
		if session.Values[SessionTypeKey] != "anonymous" {
			session.Values[SessionTypeKey] = "anonymous"

			if err := am.sess.Save(r, w, session); err != nil {
				http.Error(w, "Unable to persist anonymous session",
					http.StatusInternalServerError)
				return
			}
		} else if session.IsNew {
			if err := am.sess.Save(r, w, session); err != nil {
				http.Error(w, "Unable to persist anonymous session", http.StatusInternalServerError)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
