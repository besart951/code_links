package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/besart951/code_links/platform/auth"
	"github.com/besart951/code_links/platform/gateway"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := gateway.LoadConfig()
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

	store := gateway.NewPostgresStore(db)
	authService := auth.NewService(store, cfg.AuthConfig())
	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           gateway.NewServer(cfg, authService, store, store),
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
