package cache

import "github.com/alexbobkovv/insider-trades/pkg/redisdb"

type TradeCache struct {
	c *redisdb.RedisClient
}

func New(client *redisdb.RedisClient) *TradeCache {
	return &TradeCache{c: client}
}
