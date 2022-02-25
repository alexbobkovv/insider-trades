package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func New(url string) (*Postgres, error) {
	pool, err := pgxpool.Connect(context.Background(), url)

	if err != nil {
		return nil, fmt.Errorf("postgresql: failed to connect to db: %w", err)
	}

	return &Postgres{Pool: pool}, nil
}
