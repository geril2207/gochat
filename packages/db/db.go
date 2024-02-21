package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func DbConnect() (*pgxpool.Pool, error) {
	dbUrl := os.Getenv("DB_URL")
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return pool, err
	}
	Pool = pool

	return Pool, nil
}

func DbCloseConnection() {
	if Pool != nil {
		Pool.Close()
	}
}
