package redisdb

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	c *redis.Client
}

func New(host, port, username, password string) *RedisClient {
	redisAddr := fmt.Sprintf("%s:%s", host, port)

	options := &redis.Options{
		Addr:     redisAddr,
		Username: username,
		Password: password,
		DB:       0, // default
	}

	rdb := redis.NewClient(options)

	return &RedisClient{c: rdb}
}
