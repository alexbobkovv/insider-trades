package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func New(url string) (*Postgres, error) {
	var connAttempts = 3
	var pool *pgxpool.Pool
	const connTimeout = time.Second * 3
	var err error

	for connAttempts > 0 {
		pool, err = pgxpool.Connect(context.Background(), url)

		if err == nil {
			break
		}

		log.Printf("postgresql New: failed to connect to db, trying to reconnect.. error: %v", err)

		time.Sleep(connTimeout)
		connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgresql: failed to connect to db: %w", err)
	}

	return &Postgres{Pool: pool}, nil
}
