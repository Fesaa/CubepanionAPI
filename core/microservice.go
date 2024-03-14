package core

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type defaultMicroService[T MicroServiceConfig, D interface{}] struct {
	config T
	db     D

	app *fiber.App
}

func NewMicroService[T MicroServiceConfig, D interface{}](config T, d DatabaseProvider[D], fiberConfig ...fiber.Config) (MicroService[T, D], error) {
	if d == nil {
		return nil, fmt.Errorf("Database provider is nil")
	}

	app := fiber.New(fiberConfig...)

	db, err := d(config.Database())
	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %w", err)
	}

	m := &defaultMicroService[T, D]{
		config: config,
		app:    app,
		db:     db,
	}

	return m, nil
}

func (m *defaultMicroService[T, D]) Config() T {
	return m.config
}

func (m *defaultMicroService[T, D]) DB() D {
	return m.db
}

func (m *defaultMicroService[T, D]) App() *fiber.App {
	return m.app
}

func (m *defaultMicroService[T, D]) UseLogger(config ...logger.Config) {
	m.app.Use(logger.New(config...))
}

func (m *defaultMicroService[T, D]) UseLimiter(config ...limiter.Config) {
	m.app.Use(limiter.New(config...))
}

func (m *defaultMicroService[T, D]) UseCache(config ...cache.Config) {
	m.app.Use(cache.New(config...))
}

func (m *defaultMicroService[T, D]) UseRedisCache() {
	m.UseCache(cache.Config{
		Storage: redisCache(m.Config()),
	})
}

func (m *defaultMicroService[T, D]) Get(path string, handlers ...Handler[T, D]) fiber.Router {
	return m.app.Get(path, m.conv(handlers...)...)
}

func (m *defaultMicroService[T, D]) Post(path string, handlers ...Handler[T, D]) fiber.Router {
	return m.app.Post(path, m.conv(handlers...)...)
}

func (m *defaultMicroService[T, D]) Put(path string, handlers ...Handler[T, D]) fiber.Router {
	return m.app.Put(path, m.conv(handlers...)...)
}

func (m *defaultMicroService[T, D]) Use(args ...interface{}) fiber.Router {
	return m.app.Use(args...)
}

func (m *defaultMicroService[T, D]) UsePrometheus() {
	p := fiberprometheus.New(m.Config().ServiceName())
	p.RegisterAt(m.App(), "/metrics")
	m.App().Use(p.Middleware)
}

func (m *defaultMicroService[T, D]) UseDefaults() {
	m.UseLogger()
	m.UsePrometheus()
	m.UseLimiter(limiter.Config{
		Max:               10,
		Expiration:        10 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	})
}

func (m *defaultMicroService[T, D]) Start() error {
	slog.Info("Starting microservice", "host", m.Config().Host(), "port", m.Config().Port())
	return m.app.Listen(fmt.Sprintf("%s:%d", m.Config().Host(), m.Config().Port()))
}

func (m *defaultMicroService[T, D]) conv(handlers ...Handler[T, D]) []fiber.Handler {
	h := make([]fiber.Handler, len(handlers))
	for i, v := range handlers {
		h[i] = func(c *fiber.Ctx) error {
			return v(m, c)
		}
	}

	return h
}
