package authorisation

import (
	"context"
	"errors"
	"time"

	"booking-platform/internal/database"

	"github.com/google/uuid"
)

type Service struct {
	repo        *Repository
	transaction database.TransactorInterface
}

var (
	ErrForbidden     = errors.New("not authorized to manage role grants")
	ErrInvalidGrant  = errors.New("invalid role grant")
	ErrGrantNotFound = errors.New("role grant not found")
	ErrSystemRole    = errors.New("cannot grant this system role")
)

func NewService(repo *Repository, trx database.TransactorInterface) *Service {
	return &Service{
		repo:        repo,
		transaction: trx,
	}
}

type grantRoleInput struct {
	ActorID        uuid.UUID
	UserID         uuid.UUID
	RoleID         int32
	OrganisationID *uuid.UUID
	ExpiresAt      *time.Time
}

// HasPermissions exists to mainly be used by the Authorisation middleware
func (s *Service) HasPermissions(
	ctx context.Context,
	userID uuid.UUID,
	organisationID *uuid.UUID,
	permission PermissionKey,
) (bool, error) {
	return s.repo.HasPermissions(ctx, userID, string(permission), organisationID)
}

func (s *Service) GrantRole(
	ctx context.Context,
	input grantRoleInput,
) (*RoleGrant, error) {
	// basic validation
	if input.ActorID == uuid.Nil ||
		input.UserID == uuid.Nil ||
		input.RoleID <= 0 {
		return nil, ErrInvalidGrant
	}

	isAllowed, err := s.repo.HasPermissions(
		ctx,
		input.UserID,
		"role.grants",
		input.OrganisationID,
	)
	if err != nil {
		return nil, err
	}
	if !isAllowed {
		return nil, ErrForbidden
	}

	if input.ExpiresAt != nil &&
		input.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidGrant
	}
	err = s.transaction.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.
			CheckUserRegistration(txCtx, input.UserID, "registered"); err != nil {
			return err
		}

		if err := s.repo.PreventDuplicateActiveGrants(
			txCtx, input.UserID,
			input.RoleID, input.OrganisationID); err != nil {
			return err
		}
		grant := RoleGrant{
			UserID:         input.UserID,
			RoleID:         input.RoleID,
			OrganisationID: *input.OrganisationID,
			GrantedBy:      input.ActorID,
			ExpiresAt:      input.ExpiresAt,
		}
		return s.repo.CreateGrant(txCtx, &grant)
	})

	return nil, err
}

func (s *Service) RevokeRole(ctx context.Context, actorID, grantID uuid.UUID) error {
	if actorID == uuid.Nil || grantID == uuid.Nil {
		return ErrInvalidGrant
	}
	grant, err := s.repo.GetRoleToRevoke(ctx, actorID, grantID)
	if err != nil {
		return ErrGrantNotFound
	}
	isAllowed, err := s.repo.HasPermissions(ctx, actorID, "roles.grant", &grant.OrganisationID)
	if err != nil {
		return err
	}

	if !isAllowed {
		return ErrForbidden
	}
	return s.repo.RevokeRole(ctx, grantID)
}
