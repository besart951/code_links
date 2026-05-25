package productauth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"
)

type RemoteValidatorConfig struct {
	JWKSURL      string
	Issuer       string
	Audience     string
	ProductID    string
	Client       *http.Client
	JWKSCacheTTL time.Duration
}

type RemoteValidator struct {
	config    RemoteValidatorConfig
	mu        sync.RWMutex
	keys      map[string]*rsa.PublicKey
	expiresAt time.Time
}

type jwksDocument struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func NewRemoteValidator(ctx context.Context, config RemoteValidatorConfig) (*RemoteValidator, error) {
	if config.JWKSURL == "" {
		return nil, errors.New("JWKSURL is required")
	}
	if config.Issuer == "" {
		return nil, errors.New("Issuer is required")
	}
	if config.Audience == "" {
		return nil, errors.New("Audience is required")
	}
	if config.ProductID == "" {
		return nil, errors.New("ProductID is required")
	}
	if config.Client == nil {
		config.Client = &http.Client{Timeout: 5 * time.Second}
	}
	if config.JWKSCacheTTL <= 0 {
		config.JWKSCacheTTL = 5 * time.Minute
	}

	keys, err := fetchJWKS(ctx, config.Client, config.JWKSURL)
	if err != nil {
		return nil, err
	}

	return &RemoteValidator{config: config, keys: keys, expiresAt: time.Now().Add(config.JWKSCacheTTL)}, nil
}

func (v *RemoteValidator) refreshJWKS(ctx context.Context) error {
	keys, err := fetchJWKS(ctx, v.config.Client, v.config.JWKSURL)
	if err != nil {
		return err
	}

	v.mu.Lock()
	defer v.mu.Unlock()
	v.keys = keys
	v.expiresAt = time.Now().Add(v.config.JWKSCacheTTL)
	return nil
}

func (v *RemoteValidator) refreshJWKSIfExpired(ctx context.Context) error {
	v.mu.RLock()
	expired := time.Now().After(v.expiresAt)
	v.mu.RUnlock()
	if !expired {
		return nil
	}
	return v.refreshJWKS(ctx)
}

func (v *RemoteValidator) key(kid string) (*rsa.PublicKey, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	key, ok := v.keys[kid]
	return key, ok
}

func fetchJWKS(ctx context.Context, client *http.Client, url string) (map[string]*rsa.PublicKey, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch JWKS: unexpected status %d", response.StatusCode)
	}

	var document jwksDocument
	if err := json.NewDecoder(response.Body).Decode(&document); err != nil {
		return nil, err
	}

	keys := make(map[string]*rsa.PublicKey, len(document.Keys))
	for _, key := range document.Keys {
		publicKey, err := key.rsaPublicKey()
		if err != nil {
			return nil, fmt.Errorf("parse key %q: %w", key.Kid, err)
		}
		keys[key.Kid] = publicKey
	}

	if len(keys) == 0 {
		return nil, errors.New("JWKS did not contain any RSA keys")
	}

	return keys, nil
}

func (key jwk) rsaPublicKey() (*rsa.PublicKey, error) {
	if key.Kty != "RSA" {
		return nil, fmt.Errorf("unsupported key type %q", key.Kty)
	}
	if key.Alg != "" && key.Alg != "RS256" {
		return nil, fmt.Errorf("unsupported key alg %q", key.Alg)
	}
	if key.Use != "" && key.Use != "sig" {
		return nil, fmt.Errorf("unsupported key use %q", key.Use)
	}
	if key.Kid == "" {
		return nil, errors.New("kid is required")
	}
	if key.N == "" || key.E == "" {
		return nil, errors.New("modulus and exponent are required")
	}

	modulusBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, err
	}
	exponentBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, err
	}

	exponent := 0
	for _, b := range exponentBytes {
		exponent = exponent<<8 + int(b)
	}
	if exponent == 0 {
		return nil, errors.New("invalid exponent")
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(modulusBytes),
		E: exponent,
	}, nil
}
