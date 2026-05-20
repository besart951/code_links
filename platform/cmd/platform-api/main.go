package main

import (
	"context"
	"crypto/ed25519"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/besart951/code_links/platform/internal/access"
	adminuc "github.com/besart951/code_links/platform/internal/admin"
	authuc "github.com/besart951/code_links/platform/internal/auth"
	"github.com/besart951/code_links/platform/internal/config"
	httpapi "github.com/besart951/code_links/platform/internal/http"
	"github.com/besart951/code_links/platform/internal/passwordbcrypt"
	"github.com/besart951/code_links/platform/internal/postgres"
	"github.com/besart951/code_links/platform/internal/productclients"
	"github.com/besart951/code_links/platform/internal/shared"
	tenantuc "github.com/besart951/code_links/platform/internal/tenants"
	"github.com/besart951/code_links/platform/internal/tokenjose"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	store := postgres.NewStore(db)
	clock := shared.SystemClock{}
	ids := shared.RandomIDGenerator{}
	tokens := passwordbcrypt.RefreshTokenHasher{}

	var issueToken httpapi.IssueTokenUseCase
	tokenIssuer, err := buildTokenIssuer(cfg)
	if err != nil {
		log.Fatalf("configure token issuer: %v", err)
	}
	if tokenIssuer != nil {
		issueToken = access.IssueAccessToken{
			Repo:      store,
			Tokens:    tokenIssuer,
			Clock:     clock,
			IDs:       ids,
			Issuer:    cfg.Issuer,
			AccessTTL: cfg.AccessTokenTTL,
		}
	}

	handler := httpapi.NewServer(httpapi.Server{
		Login: authuc.LoginUser{
			Users:           store,
			Sessions:        store,
			Passwords:       passwordbcrypt.Verifier{},
			Tokens:          tokens,
			Clock:           clock,
			IDs:             ids,
			SessionTTL:      cfg.SessionTTL,
			RefreshTokenTTL: cfg.RefreshTokenTTL,
		},
		Refresh: authuc.RefreshSession{
			Users:           store,
			Sessions:        store,
			Tokens:          tokens,
			Clock:           clock,
			IDs:             ids,
			RefreshTokenTTL: cfg.RefreshTokenTTL,
		},
		Logout: authuc.LogoutSession{
			Sessions: store,
			Clock:    clock,
		},
		Authenticate: authuc.AuthenticateSession{
			Sessions: store,
			Tokens:   tokens,
			Clock:    clock,
		},
		SelectTenant: tenantuc.SelectTenant{Tenants: store},
		IssueToken:   issueToken,
		CheckProduct: access.CheckProductAccess{Repo: store, Clock: clock},
		CheckFeature: access.CheckFeatureAccess{Repo: store, Clock: clock},
		ProductClients: productclients.StaticAuthenticator{
			Tokens: cfg.ProductClientTokens,
		},
		Tenants:      store,
		Products:     store,
		Entitlements: store,
		Admin: adminuc.ReadService{
			Repo:  store,
			Clock: clock,
		},
		Config: httpapi.Config{
			CookieDomain: cfg.CookieDomain,
			CookieSecure: cfg.CookieSecure,
		},
	})

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("platform-api listening on %s", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("serve platform-api: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown platform-api: %v", err)
	}
}

func buildTokenIssuer(cfg config.Config) (access.TokenIssuer, error) {
	if cfg.SigningPrivateKeyPath == "" {
		log.Print("TOKEN_SIGNING_PRIVATE_KEY_PATH missing; /api/v1/auth/token disabled")
		return nil, nil
	}
	if len(cfg.AudiencePublicKeyPaths) == 0 {
		log.Print("AUDIENCE_PUBLIC_KEY_PATHS missing; /api/v1/auth/token disabled")
		return nil, nil
	}
	privateKey, err := tokenjose.LoadEd25519PrivateKey(cfg.SigningPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	registry := tokenjose.StaticRegistry{
		Signing: tokenjose.SigningKey{
			KeyID:      cfg.SigningKeyID,
			PrivateKey: privateKey,
			PublicKey:  privateKey.Public().(ed25519.PublicKey),
		},
		Encryption: map[access.ProductKey]tokenjose.AudienceEncryptionKey{},
		Decryption: map[access.ProductKey]tokenjose.AudienceDecryptionKey{},
	}
	for audience, path := range cfg.AudiencePublicKeyPaths {
		publicKey, err := tokenjose.LoadECDSAPublicKey(path)
		if err != nil {
			return nil, err
		}
		registry.Encryption[access.ProductKey(audience)] = tokenjose.AudienceEncryptionKey{
			Audience:  access.ProductKey(audience),
			KeyID:     audience + "-enc-1",
			PublicKey: publicKey,
		}
	}
	return tokenjose.Issuer{Keys: registry}, nil
}
