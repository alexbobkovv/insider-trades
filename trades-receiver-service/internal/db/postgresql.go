package postgresql

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func New(url string, poolMax int) (*Postgres, error) {
	pool, err := pgxpool.Connect(context.Background(), url)

	if err != nil {
		return nil, err
	}

	return &Postgres{Pool: pool}, nil
}