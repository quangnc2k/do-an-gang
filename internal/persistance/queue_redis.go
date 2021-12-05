package persistance

import (
	"context"
	"github.com/go-redis/redis"
	"time"
)

type RedisQueue struct {
	client *redis.Client
}

func (q *RedisQueue) Pop(ctx context.Context) (channel string, data string, err error) {
	val, err := q.client.BLPop(30*time.Second, "files", "notice", "conn").Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			err = ErrOutOfItem
			return
		}
	}
	if len(val) != 2 {
		return
	}

	return val[0], val[1], err
}

func NewRedisClient(url string, poolSize int) (*redis.Client, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	options.PoolSize = poolSize
	options.IdleTimeout = -1

	redisClient := redis.NewClient(options)
	return redisClient, nil
}

var _ QueueStore = (*RedisQueue)(nil)
