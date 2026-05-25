package productauth

import "context"

type claimsContextKey struct{}
type currentUserContextKey struct{}

func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey{}).(Claims)
	return claims, ok
}

func contextWithClaims(ctx context.Context, claims Claims) context.Context {
	ctx = context.WithValue(ctx, claimsContextKey{}, claims)
	return context.WithValue(ctx, currentUserContextKey{}, claims.CurrentUser())
}

func CurrentUserFromContext(ctx context.Context) (CurrentUser, bool) {
	user, ok := ctx.Value(currentUserContextKey{}).(CurrentUser)
	return user, ok
}
