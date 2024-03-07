package core

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type defaultMicroService[T MicroServiceConfig] struct {
	config T
	db     *sql.DB

	app *fiber.App
}

func NewMicroService[T MicroServiceConfig](config T, d DatabaseProvider, fiberConfig ...fiber.Config) (MicroService[T], error) {
	if d == nil {
		return nil, fmt.Errorf("Database provider is nil")
	}

	app := fiber.New(fiberConfig...)

	db, err := d(config.Database())
	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %w", err)
	}

	m := &defaultMicroService[T]{
		config: config,
		app:    app,
		db:     db,
	}

	return m, nil
}

func (m *defaultMicroService[T]) Config() T {
	return m.config
}

func (m *defaultMicroService[T]) DB() *sql.DB {
	return m.db
}

func (m *defaultMicroService[T]) App() *fiber.App {
	return m.app
}

func (m *defaultMicroService[T]) UseLogger(config ...logger.Config) {
	m.app.Use(logger.New(config...))
}

func (m *defaultMicroService[T]) UseLimiter(config ...limiter.Config) {
	m.app.Use(limiter.New(config...))
}

func (m *defaultMicroService[T]) UseCache(config ...cache.Config) {
	m.app.Use(cache.New(config...))
}

func (m *defaultMicroService[T]) UseRedisCache() {
	m.UseCache(cache.Config{
		Storage: redisCache(m.Config()),
	})
}

func (m *defaultMicroService[T]) Get(path string, handlers ...Handler[T]) fiber.Router {
	return m.app.Get(path, m.conv(handlers...)...)
}

func (m *defaultMicroService[T]) Post(path string, handlers ...Handler[T]) fiber.Router {
	return m.app.Post(path, m.conv(handlers...)...)
}

func (m *defaultMicroService[T]) Put(path string, handlers ...Handler[T]) fiber.Router {
	return m.app.Put(path, m.conv(handlers...)...)
}

func (m *defaultMicroService[T]) Use(args ...interface{}) fiber.Router {
	return m.app.Use(args...)
}

func (m *defaultMicroService[T]) UsePrometheus() {
	p := fiberprometheus.New(m.Config().ServiceName())
	p.RegisterAt(m.App(), "/metrics")
}

func (m *defaultMicroService[T]) UseDefaults() {
	m.UseLogger()
	m.UsePrometheus()
	m.UseLimiter(limiter.Config{
		Max:               10,
		Expiration:        10 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	})
}

func (m *defaultMicroService[T]) Start() error {
	slog.Info("Starting microservice", "host", m.Config().Host(), "port", m.Config().Port())
	return m.app.Listen(fmt.Sprintf("%s:%d", m.Config().Host(), m.Config().Port()))
}

func (m *defaultMicroService[T]) conv(handlers ...Handler[T]) []fiber.Handler {
	h := make([]fiber.Handler, len(handlers))
	for i, v := range handlers {
		h[i] = func(c *fiber.Ctx) error {
			return v(m, c)
		}
	}

	return h
}
