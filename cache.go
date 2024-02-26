package main

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

func redisCache(c RedisConfig) fiber.Storage {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Address,
		Password: c.PassWord,
		DB:       c.DB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Info("Cannot connect to redis, falling back to in-memory cache")
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
