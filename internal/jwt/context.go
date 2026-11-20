package jwt

import (
	"context"

	claims "booking-platform/internal/claims"
)

type contextKey string

const claimContextKey contextKey = "authenticaed_user"

func WithClaims(
	ctx context.Context, userclaims *claims.Claims,
) context.Context {
	return context.WithValue(ctx, claimContextKey, userclaims)
}

func ClaimsFromContext(ctx context.Context) (*claims.Claims, bool) {
	userclaims, ok := ctx.Value(claimContextKey).(*claims.Claims)
	return userclaims, ok
}
