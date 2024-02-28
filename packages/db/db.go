package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DbConnect(dbUrl string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return pool, err
	}

	return pool, nil
}

func DbCloseConnection(pool *pgxpool.Pool) error {
	if pool != nil {
		pool.Close()
	}
	return nil
}
