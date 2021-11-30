package persistance

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/config"
	"log"

	"github.com/pkg/errors"
)

var queue QueueStore

type QueueStore interface {
	Pop(ctx context.Context) (channel string, data string, err error)
}

func GetQueue() QueueStore {
	return queue
}

var ErrOutOfItem = errors.New("queue empty")

func LoadQueueWithRedis(ctx context.Context) {
	rdb, err := NewRedisClient(config.Env.RedisURL, 50)
	if err != nil {
		log.Fatalln("cannot init redis queue:", err)
	}

	q := RedisQueue{client: rdb}

	queue = &q
}