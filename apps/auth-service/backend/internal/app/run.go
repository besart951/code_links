package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/besart951/code-links/apps/auth-service/backend/internal/admin"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/auth"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/config"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/httpapi"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/mail"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/store/memory"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/store/postgres"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/token"
)

type store interface {
	auth.Store
	auth.RoleStore
	auth.NotificationStore
	admin.Store
	admin.SettingsStore
	admin.NotificationStore
	admin.AuditStore
	admin.SessionStore
}

func Run(ctx context.Context) error {
	cfg := config.Load()

	signer, err := token.NewSigner(token.Config{
		KeyID:          cfg.JWTKeyID,
		Issuer:         cfg.Issuer,
		Audience:       cfg.Audience,
		Lifetime:       cfg.AccessTokenLifetime,
		PrivateKeyPEM:  cfg.JWTPrivateKeyPEM,
		PrivateKeyFile: cfg.JWTPrivateKeyFile,
	})
	if err != nil {
		return err
	}

	store, cleanup, err := openStore(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer cleanup()

	authService := auth.NewService(auth.Config{
		Environment:          cfg.Environment,
		PublicFrontendURL:    cfg.PublicFrontendURL,
		RefreshTokenLifetime: cfg.RefreshTokenLifetime,
	}, store, store, store, signer)
	adminService := admin.NewService(admin.Config{
		SMTPSecretKey: cfg.SMTPSecretKey,
	}, store, store, store, store, store, signer, mail.SmtpSender{})
	server := httpapi.NewServer(httpapi.Config{
		AllowedOrigins:     cfg.AllowedOrigins,
		EnableMockPurchase: cfg.EnableMockPurchase,
		CookieDomain:       cfg.CookieDomain,
		CookieSecure:       cfg.CookieSecure,
		PublicFrontendURL:  cfg.PublicFrontendURL,
	}, authService, adminService, signer)

	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           server.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("auth-service listening on :%s", cfg.Port)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func openStore(ctx context.Context, databaseURL string) (store, func(), error) {
	if databaseURL == "" {
		store, err := memory.New()
		return store, func() {}, err
	}
	return postgres.Open(ctx, databaseURL)
}
