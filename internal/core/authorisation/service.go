package authorisation

import (
	"booking-platform/internal/database"
	"context"
	"errors"
	"time"

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
	RoleId         int32
	OrganisationId *uuid.UUID
	ExpiresAt      *time.Time
}

func (s *Service) GrantRole(
	ctx context.Context,
	input grantRoleInput) (*RoleGrant, error) {
	//basic validation
	if input.ActorID == uuid.Nil ||
		input.UserID == uuid.Nil ||
		input.RoleId <= 0 {
		return nil, ErrInvalidGrant
	}

	isAllowed, err := s.repo.HasPermissions(
		ctx,
		input.UserID,
		"role.grants",
		input.OrganisationId,
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
			input.RoleId, input.OrganisationId); err != nil {
			return err
		}
		grant := RoleGrant{
			UserId:         input.UserID,
			RoleId:         input.RoleId,
			OrganisationId: *input.OrganisationId,
			GrantedBy:      input.ActorID,
			ExpiresAt:      input.ExpiresAt,
		}
		return s.repo.CreateGrant(txCtx, &grant)
	})

	return nil, nil
}
func (s *Service) RevokeRole(ctx context.Context, actorId, grantId uuid.UUID) error {
	if actorId == uuid.Nil || grantId == uuid.Nil {
		return ErrInvalidGrant
	}
	grant, err := s.repo.GetRoleToRevoke(ctx, actorId, grantId)
	if err != nil {
		return ErrGrantNotFound
	}
	isAllowed, err := s.repo.HasPermissions(ctx, actorId, "roles.grant", &grant.OrganisationId)

	if err != nil {
		return err
	}

	if !isAllowed {
		return ErrForbidden
	}
	return s.repo.RevokeRole(ctx, grantId)
}
