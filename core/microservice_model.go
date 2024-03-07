package core

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Handler[T MicroServiceConfig] func(MicroService[T], *fiber.Ctx) error

type MicroService[T MicroServiceConfig] interface {
	Config() T
	DB() *sql.DB

	App() *fiber.App
	Get(string, ...Handler[T]) fiber.Router
	Post(string, ...Handler[T]) fiber.Router
	Put(string, ...Handler[T]) fiber.Router
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
