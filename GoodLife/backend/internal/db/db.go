package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect returns a connection pool to the PostgreSQL database.
func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}
	// Verify the connection
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}
