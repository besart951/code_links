package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/besart951/code_links/platform/auth"
	"github.com/besart951/code_links/platform/gateway"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	email := strings.TrimSpace(os.Getenv("SEED_USER_EMAIL"))
	password := os.Getenv("SEED_USER_PASSWORD")
	if email == "" || password == "" {
		log.Print("SEED_USER_EMAIL or SEED_USER_PASSWORD missing; skipping platform seed")
		return
	}

	cfg := gateway.LoadConfig()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	if err := seed(ctx, db, email, password); err != nil {
		log.Fatal(err)
	}
}

func seed(ctx context.Context, db *pgxpool.Pool, email, password string) error {
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	displayName := env("SEED_USER_DISPLAY_NAME", "CodeLinks Admin")
	tenantName := env("SEED_TENANT_NAME", "CodeLinks")
	tenantSlug := env("SEED_TENANT_SLUG", "codelinks")

	var userID string
	if err := db.QueryRow(ctx, `
		insert into users (email, password_hash, display_name, status)
		values (lower($1), $2, $3, 'active')
		on conflict (email) do update set display_name = excluded.display_name
		returning id::text
	`, email, passwordHash, displayName).Scan(&userID); err != nil {
		return err
	}

	var roleID string
	if err := db.QueryRow(ctx, `select id::text from roles where key = 'owner'`).Scan(&roleID); err != nil {
		return err
	}

	var tenantID string
	if err := db.QueryRow(ctx, `
		insert into tenants (type, name, slug, owner_user_id, status, billing_email)
		values ('company', $1, $2, $3::uuid, 'active', $4)
		on conflict (slug) do update set name = excluded.name
		returning id::text
	`, tenantName, tenantSlug, userID, email).Scan(&tenantID); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `
		insert into tenant_members (tenant_id, user_id, role_id, status)
		values ($1::uuid, $2::uuid, $3::uuid, 'active')
		on conflict (tenant_id, user_id) do update set role_id = excluded.role_id, status = 'active'
	`, tenantID, userID, roleID); err != nil {
		return err
	}

	for _, feature := range seedEntitlements() {
		parts := strings.SplitN(feature, ":", 2)
		if len(parts) != 2 {
			continue
		}
		if _, err := db.Exec(ctx, `
			insert into entitlements (tenant_id, product_key, feature_key, source, enabled)
			values ($1::uuid, $2, $3, 'manual', true)
			on conflict (tenant_id, product_key, feature_key, source) do update set enabled = true
		`, tenantID, parts[0], parts[1]); err != nil {
			return err
		}
	}

	log.Printf("seeded platform user %s and tenant %s", email, tenantSlug)
	return nil
}

func seedEntitlements() []string {
	raw := env("SEED_ENTITLEMENTS", strings.Join([]string{
		"planer_link:planer.pdf_export",
		"planer_link:planer.excel_export",
		"infra_link:infra.module_bacnet",
		"infra_link:infra.module_sps",
		"infra_link:infra.module_field_devices",
	}, ","))
	return strings.Split(raw, ",")
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
