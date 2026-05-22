package productauth

import (
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

	return v.ValidateToken(token)
}

func (v *RemoteValidator) ValidateToken(rawToken string) (Claims, error) {
	claims := Claims{}
	parser := jwt.NewParser(
		jwt.WithIssuer(v.config.Issuer),
		jwt.WithAudience(v.config.Audience),
		jwt.WithExpirationRequired(),
	)

	token, err := parser.ParseWithClaims(rawToken, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}

		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, errors.New("missing kid")
		}

		key, ok := v.keys[kid]
		if !ok {
			return nil, errors.New("unknown kid")
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
