package postgres

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/besart951/code_links/platform/internal/access"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestStoreQueriesAgainstMigratedDatabase(t *testing.T) {
	dsn := os.Getenv("PLATFORM_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("PLATFORM_TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store := NewStore(db)
	if _, err := store.ListProducts(ctx); err != nil {
		t.Fatalf("list products query failed: %v", err)
	}
	if _, err := store.SessionAccess(ctx, access.SessionID("missing_session")); !errors.Is(err, access.ErrSessionInactive) {
		t.Fatalf("expected missing session to map to ErrSessionInactive, got %v", err)
	}
}
