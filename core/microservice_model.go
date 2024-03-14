package core

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Handler[T MicroServiceConfig, D interface{}] func(MicroService[T, D], *fiber.Ctx) error

type MicroService[T MicroServiceConfig, D interface{}] interface {
	Config() T
	DB() D

	App() *fiber.App
	Get(string, ...Handler[T, D]) fiber.Router
	Post(string, ...Handler[T, D]) fiber.Router
	Put(string, ...Handler[T, D]) fiber.Router
	Use(args ...interface{}) fiber.Router

	// This registers UsePrometheus, UseLimiter, and UseLogger with sensible defaults
	UseDefaults()

	// Registers the fiberprometheus middleware with MicroServiceConfig.ServiceName as the namespace
	UsePrometheus()
	UseLogger(...logger.Config)
	UseLimiter(...limiter.Config)
	UseCache(...cache.Config)

	// Registers the redis cache, will automatically fall back to in-memory cache if redis is not available
	UseRedisCache()

	Start() error
}
