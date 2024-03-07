package core

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisWrapper struct {
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

	return &RedisWrapper{rdb}
}

func (rw *RedisWrapper) Get(key string) ([]byte, error) {
	return rw.rdb.Get(ctx, key).Bytes()
}

func (rw *RedisWrapper) Set(key string, val []byte, exp time.Duration) error {
	return rw.rdb.Set(ctx, key, val, exp).Err()
}

func (rw *RedisWrapper) Delete(key string) error {
	return rw.rdb.Del(ctx, key).Err()
}

func (rw *RedisWrapper) Reset() error {
	return rw.rdb.FlushDB(ctx).Err()
}

func (rw *RedisWrapper) Close() error {
	return rw.rdb.Close()
}
