package redisdb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	*redis.Client
}

func New(host, port, password string) (*RedisClient, error) {
	redisAddr := fmt.Sprintf("%s:%s", host, port)

	options := &redis.Options{
		Addr:     redisAddr,
		Password: password,
		DB:       0, // default
	}

	rdb := redis.NewClient(options)

	var connAttempts = 5
	const connTimeout = time.Second * 3
	var err error

	for connAttempts > 0 {
		status := rdb.Ping(context.Background())
		err = status.Err()

		if err == nil {
			break
		}

		log.Printf("redis New: failed to ping redis, retrying.. error: %v", err)

		time.Sleep(connTimeout)
		connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("redis New: failed to connect to redis: %w", err)
	}

	return &RedisClient{rdb}, nil
}
