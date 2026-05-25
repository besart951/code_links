package productauth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (v *RemoteValidator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := v.ValidateRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !claims.HasLicense(v.config.ProductID) {
			http.Error(w, "missing product license", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r.WithContext(contextWithClaims(r.Context(), claims)))
	})
}

func (v *RemoteValidator) ValidateRequest(r *http.Request) (Claims, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return Claims{}, errors.New("missing Authorization header")
	}

	token, ok := strings.CutPrefix(header, "Bearer ")
	if !ok || token == "" {
		return Claims{}, errors.New("Authorization header must be Bearer token")
	}

	return v.ValidateTokenWithContext(r.Context(), token)
}

func (v *RemoteValidator) ValidateToken(rawToken string) (Claims, error) {
	return v.ValidateTokenWithContext(context.Background(), rawToken)
}

func (v *RemoteValidator) ValidateTokenWithContext(ctx context.Context, rawToken string) (Claims, error) {
	if err := v.refreshJWKSIfExpired(ctx); err != nil {
		return Claims{}, err
	}

	claims := Claims{}
	refreshedUnknownKid := false
	parser := jwt.NewParser(
		jwt.WithIssuer(v.config.Issuer),
		jwt.WithAudience(v.config.Audience),
		jwt.WithExpirationRequired(),
	)

	token, err := parser.ParseWithClaims(rawToken, &claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodRS256 {
			return nil, errors.New("unexpected signing method")
		}

		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, errors.New("missing kid")
		}

		key, ok := v.key(kid)
		if !ok {
			if refreshedUnknownKid {
				return nil, errors.New("unknown kid")
			}
			refreshedUnknownKid = true
			if err := v.refreshJWKS(ctx); err != nil {
				return nil, err
			}
			key, ok = v.key(kid)
			if !ok {
				return nil, errors.New("unknown kid")
			}
		}

		return key, nil
	})
	if err != nil {
		return Claims{}, err
	}
	if !token.Valid {
		return Claims{}, errors.New("invalid token")
	}

	return claims, nil
}
