package postgres

import (
	"context"
	"embed"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, `
		create table if not exists schema_migrations (
			version text primary key,
			applied_at timestamptz not null default now()
		)
	`); err != nil {
		return err
	}

	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return err
	}
	versions := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			versions = append(versions, entry.Name())
		}
	}
	sort.Strings(versions)

	for _, version := range versions {
		applied, err := migrationApplied(ctx, pool, version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		content, err := migrationFiles.ReadFile("migrations/" + version)
		if err != nil {
			return err
		}
		if err := applyMigration(ctx, pool, version, string(content)); err != nil {
			return err
		}
	}

	return nil
}

func migrationApplied(ctx context.Context, pool *pgxpool.Pool, version string) (bool, error) {
	var applied bool
	err := pool.QueryRow(ctx, `select exists(select 1 from schema_migrations where version = $1)`, version).Scan(&applied)
	return applied, err
}

func applyMigration(ctx context.Context, pool *pgxpool.Pool, version string, sql string) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if _, err := conn.Conn().PgConn().Exec(ctx, sql).ReadAll(); err != nil {
		return err
	}
	_, err = conn.Exec(ctx, `insert into schema_migrations (version) values ($1)`, version)
	return err
}
