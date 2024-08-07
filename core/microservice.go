package core

import (
	"fmt"
	"github.com/Fesaa/CubepanionAPI/core/log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
)

type defaultMicroService[T MicroServiceConfig, D interface{}] struct {
	config T
	db     D

	app *fiber.App
}

func NewMicroService[T MicroServiceConfig, D Database](config T, d DatabaseProvider[D], fiberConfig ...fiber.Config) (MicroService[T, D], error) {
	if d == nil {
		return nil, fmt.Errorf("Database provider is nil")
	}

	app := fiber.New(fiberConfig...)

	db, err := d(config.Database())
	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %w", err)
	}

	opt := &slog.HandlerOptions{
		AddSource:   config.LoggingConfig().Source(),
		Level:       config.LoggingConfig().LogLevel(),
		ReplaceAttr: nil,
	}
	var h slog.Handler
	switch strings.ToUpper(config.LoggingConfig().Handler()) {
	case "TEXT":
		h = slog.NewTextHandler(os.Stdout, opt)
	case "JSON":
		h = slog.NewJSONHandler(os.Stdout, opt)
	default:
		panic("Invalid logging handler: " + config.LoggingConfig().Handler())
	}
	_log := slog.New(h)
	slog.SetDefault(_log)
	log.SetDefault(_log)

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
	var c logger.Config
	if len(config) > 0 {
		c = config[0]
	} else {
		c = logger.Config{}
	}

	if c.Next == nil {
		c.Next = func(c *fiber.Ctx) bool {
			return !m.Config().LoggingConfig().LogHttp()
		}
	}

	m.app.Use(logger.New(c))
}

func (m *defaultMicroService[T, D]) UseLimiter(config ...limiter.Config) {
	m.app.Use(limiter.New(config...))
}

func (m *defaultMicroService[T, D]) UseCache(config ...cache.Config) {
	var c cache.Config
	if len(config) > 0 {
		c = config[0]
	} else {
		c = cache.Config{}
	}

	if c.KeyGenerator == nil {
		c.KeyGenerator = func(c *fiber.Ctx) string {
			return m.Config().ServiceName() + "_" + utils.CopyString(c.Context().URI().String())
		}
	}

	m.app.Use(cache.New(c))
}

func (m *defaultMicroService[T, D]) UseRedisCache(config ...cache.Config) {
	var c cache.Config
	if len(config) > 0 {
		c = config[0]
	} else {
		c = cache.Config{}
	}

	c.Storage = redisCache(m.Config())
	m.UseCache(c)
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
	m.UsePrometheus()
	m.UseLogger()
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
