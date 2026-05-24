package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

func Open(ctx context.Context, databaseURL string) (*Store, func(), error) {
	var lastErr error
	for attempt := 0; attempt < 30; attempt++ {
		pool, err := pgxpool.New(ctx, databaseURL)
		if err == nil {
			if err = pool.Ping(ctx); err == nil {
				if err = runMigrations(ctx, pool); err != nil {
					pool.Close()
					return nil, func() {}, err
				}
				return &Store{pool: pool}, pool.Close, nil
			}
			pool.Close()
		}
		lastErr = err
		time.Sleep(time.Second)
	}

	return nil, func() {}, lastErr
}
