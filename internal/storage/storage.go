package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func GetConnect(connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
