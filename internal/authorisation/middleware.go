package authorisation

import (
	"booking-platform/internal/jwt"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AuthorisationMiddleware struct {
	svc *Service
}

func NewAuthorisationMiddleware(service *Service) *AuthorisationMiddleware {
	return &AuthorisationMiddleware{
		svc: service,
	}
}

func (m *AuthorisationMiddleware) RequirePermission(
	permission PermissionKey,
	/*implies we're getting the org ID from the request*/
	organisationID func(*http.Request) (*uuid.UUID, error),
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userClaim, ok := jwt.ClaimsFromContext(r.Context())
			if !ok {
				http.Error(w, "unauthorised", http.StatusUnauthorized)
				return
			}

			orgID, err := organisationID(r)
			if err != nil {
				http.Error(w, "invalid organisation", http.StatusBadRequest)
				return
			}

			allowed, err := m.svc.HasPermissions(r.Context(), userClaim.UserID, orgID, permission)
			if err != nil {
				http.Error(w, "unable to verify permission", http.StatusInternalServerError)
				return
			}

			if !allowed {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func OrganisationIDFromURL(r *http.Request) (*uuid.UUID, error) {
	id := chi.URLParam(r, "organisationId")
	if id == "" {
		return nil, fmt.Errorf("organisation Id is required")
	}
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
