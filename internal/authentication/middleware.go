package authentication

import (
	"booking-platform/internal/core"
	"booking-platform/internal/core/authorisation"
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

const (
	sessionTypeKey = "app_session"
	jwtTypeKey     = "jwt_token"
)

type authenticatedUser struct {
	id   uuid.UUID
	role core.Role
}
type permissionsInterface interface {
	HasPermissions(context.Context, uuid.UUID, string, *uuid.UUID) (bool, error)
}

type jwtTokenValidator interface {
	ValidateToken(tokenString string) (*claims, error)
	CreateNewToken(userId uuid.UUID) (string, error)
}

type AppMiddleware struct {
	db    *gorm.DB
	name  string
	sess  sessions.Store
	store *storeManager
	jwt   jwtTokenValidator
	perm  permissionsInterface
}

func NewAppMiddleware(store *storeManager, name string, jwtInterface jwtTokenValidator, permissions permissionsInterface) *AppMiddleware {
	return &AppMiddleware{
		sess:  store.session,
		store: store,
		name:  name,
		jwt:   jwtInterface,
		perm:  permissions,
	}
}

// godoc
// Authentication Layer using jwt Tokens from claims
func (am *AppMiddleware) AuthMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//this token comes from the session
		tokenStr, err := am.store.GetToken(r)
		if err != nil {
			http.Error(w, "Cannot Create A Provisioned Guest", http.StatusInternalServerError)
			return
		}

		claims, err := am.jwt.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorised: No Valid Token Or It Expired", http.StatusUnauthorized)
			return
		}

		// authUser := authenticatedUser{
		// 	id:   claims.UserId,
		// 	role: claims.Role,
		// }
		//find the one record from the claims
		var user struct {
			ID            uuid.UUID
			AccountStatus string
		}
		if err := am.db.WithContext(r.Context()).Table("users").
			Select("id, account_status").
			Where("id = ? AND deleted_at IS NULL ", claims.ID).
			Take(&user).Error; err != nil {
			http.Error(w, "Unauthorised Login Account", http.StatusUnauthorized)
			return
		}

		if user.AccountStatus != "registered" {
			http.Error(w, "Unauthorised Login To Unregistered Account", http.StatusUnauthorized)
			return
		}

		principal := authorisation.Principal{
			UserID: claims.UserId,
			Role:   claims.Role,
		}

		ctx := authorisation.WithPrincipal(r.Context(), principal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// godoc
// Ensures every visitor (especially anonymouse) has a temp database record and a valid session token
func (am *AppMiddleware) EnsureAnonymous(next http.Handler) http.Handler {
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

// Expects the organisation scope to be provided by trusted route context or
// resolved from the requested resource.
// Do not blindly trust an organization ID from a query parameter.
// Resolve the organization from the target resource or validate the
// caller's access to it before using it as the authorization scope.
// For example, a booking-management endpoint should load the booking,
// determine its organization, and then check bookings.manage against that
// organization.
func (am *AppMiddleware) RequirePermission(
	permissionKey string,
	organisationFromRequest func(*http.Request) (*uuid.UUID, error),
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, ok := authorisation.PrincipalFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorised", http.StatusUnauthorized)
				return
			}
			var orgId *uuid.UUID
			var err error
			if organisationFromRequest != nil {
				orgId, err = organisationFromRequest(r)
				if err != nil {
					http.Error(w, "Invalid Organisation Scope", http.StatusBadRequest)
					return
				}
			}
			isAllowed, err := am.perm.HasPermissions(
				r.Context(), principal.UserID, permissionKey, orgId)
			if err != nil {
				http.Error(w, "authorisation unavailable", http.StatusInternalServerError)
				return
			}
			if !isAllowed {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func AnonymousTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(anonymousKey).(string)
	return token, ok && token != ""
}
