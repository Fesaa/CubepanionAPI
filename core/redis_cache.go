package core

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type redisWrapper struct {
	rdb *redis.Client
}

func redisCache(c MicroServiceConfig) fiber.Storage {

	rdb := redis.NewClient(&redis.Options{
		Addr:           c.Redis().Host(),
		Password:       c.Redis().Password(),
		DB:             c.Redis().DB(),
		ClientName:     c.ServiceName(),
		IdentitySuffix: "cubepanion-",
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Error("Cannot connect to redis, falling back to in-memory cache", "error", err)
		return nil
	}

	return &redisWrapper{rdb}
}

func (rw *redisWrapper) Get(key string) ([]byte, error) {
	return rw.rdb.Get(ctx, key).Bytes()
}

func (rw *redisWrapper) Set(key string, val []byte, exp time.Duration) error {
	return rw.rdb.Set(ctx, key, val, exp).Err()
}

func (rw *redisWrapper) Delete(key string) error {
	return rw.rdb.Del(ctx, key).Err()
}

func (rw *redisWrapper) Reset() error {
	return rw.rdb.FlushDB(ctx).Err()
}

func (rw *redisWrapper) Close() error {
	return rw.rdb.Close()
}
