package productclients

import (
	"context"
	"crypto/subtle"
	"errors"

	"github.com/besart951/code_links/platform/internal/access"
)

var ErrInvalidProductClient = errors.New("invalid product client")

type ProductClient struct {
	ID       string
	Audience access.ProductKey
}

type StaticAuthenticator struct {
	Tokens map[string]string
}

func (a StaticAuthenticator) Authenticate(ctx context.Context, productKey, bearerToken string) (ProductClient, error) {
	expected := a.Tokens[productKey]
	if expected == "" || bearerToken == "" {
		return ProductClient{}, ErrInvalidProductClient
	}
	if subtle.ConstantTimeCompare([]byte(expected), []byte(bearerToken)) != 1 {
		return ProductClient{}, ErrInvalidProductClient
	}
	return ProductClient{ID: productKey, Audience: access.ProductKey(productKey)}, nil
}
