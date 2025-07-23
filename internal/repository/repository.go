package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func (r *Repository) Close() {
	r.pool.Close()
}

func New(DSN string) (*Repository, error) {
	// Parse the DSN and create a new connection pool.
	cfg, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		return nil, err
	}
	// Create a new connection pool with the parsed configuration.
	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	// Ping the database to ensure the connection is valid.
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	repo := &Repository{
		pool: pool,
	}

	return repo, nil
}
