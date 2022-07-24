package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/alexbobkovv/insider-trades/api"
	"github.com/alexbobkovv/insider-trades/pkg/logger"
	"github.com/alexbobkovv/insider-trades/pkg/redisdb"
	"github.com/alexbobkovv/insider-trades/pkg/types/cursor"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)

type TradeCache struct {
	client *redisdb.RedisClient
	l      *logger.Logger
}

func New(client *redisdb.RedisClient, logger *logger.Logger) *TradeCache {
	return &TradeCache{client: client, l: logger}
}

func (c *TradeCache) ListTrades(ctx context.Context, reqCursor *cursor.Cursor, limit uint32) ([]*api.TradeViewResponse, *cursor.Cursor, error) {
	const methodName = "(c *TradeCache) ListTradesJSON"
	const setName = "tradeViews"

	var redisCursor string
	if reqCursor.IsEmpty() {
		redisCursor = "+inf"
	} else {
		cursorUNIX := reqCursor.GetUNIXTime()
		redisCursor = strconv.Itoa(int(cursorUNIX))
	}

	// ZRANGEBYSCORE zset -inf +inf WITHSCORES LIMIT 0 2
	tradeStrings, err := c.client.ZRevRangeByScore(ctx, setName, &redis.ZRangeBy{
		Min:    "0",
		Max:    redisCursor,
		Offset: 0,
		Count:  int64(limit),
	}).Result()

	if err != nil {
		c.l.Errorf("%s: failed to scan trades range from cache: %w", methodName, err)
		return nil, nil, fmt.Errorf("%s: failed to scan trades range from cache: %w", methodName, err)
	}

	if len(tradeStrings) == 0 || len(tradeStrings) != int(limit) {
		return nil, nil, nil
	}

	tradeViews := make([]*api.TradeViewResponse, len(tradeStrings))

	for idx, tradeString := range tradeStrings {
		view := &api.TradeViewResponse{}
		if err := proto.Unmarshal([]byte(tradeString), view); err != nil {
			return nil, nil, fmt.Errorf("%s: failed to unmarshal redis tradeView to *api.TradeViewResponse: %w", methodName, err)
		}

		tradeViews[idx] = view
	}

	lastView := tradeViews[len(tradeViews)-1]
	var nextCursor *cursor.Cursor

	if lastView.CreatedAt != nil {
		createdAtTime := lastView.CreatedAt.AsTime()
		nextCursor = cursor.NewFromTime(&createdAtTime)
	}

	return tradeViews, nextCursor, nil
}

func (c *TradeCache) AddTrades(ctx context.Context, trades []*api.TradeViewResponse) {
	const methodName = "(c *TradeCache) AddTrades"
	const setName = "tradeViews"

	tradesRedisZ := make([]*redis.Z, 0, len(trades))

	for _, trade := range trades {
		if trade == nil || trade.CreatedAt == nil {
			c.l.Errorf("%s: got empty trade or field: %v", methodName, trade)
			continue
		}
		createdAt := trade.CreatedAt.AsTime()
		cursorTime := cursor.NewFromTime(&createdAt)
		cursorUNIX := cursorTime.GetUNIXTime()

		tradeBytes, err := proto.Marshal(trade)
		if err != nil {
			c.l.Errorf("%s: failed to marshal trade to bytes %w", methodName, err)
			return
		}

		tradesRedisZ = append(tradesRedisZ,
			&redis.Z{
				// TODO fix in 2254, float64 UNIX time expires in Tue Jun 05 2255
				Score:  float64(cursorUNIX),
				Member: tradeBytes,
			})

	}

	cmd := c.client.ZAddNX(ctx, setName, tradesRedisZ...)
	err := cmd.Err()
	if err != nil {
		c.l.Errorf("%s: failed to add trades to cache: %s", methodName, err)
		return
	}
}
