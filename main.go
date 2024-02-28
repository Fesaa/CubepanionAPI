package main

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Fesaa/CubepanionAPI/impl"
	"github.com/Fesaa/CubepanionAPI/integration"
	"github.com/Fesaa/CubepanionAPI/routes"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	config, err := NewConfig("config.yaml")
	if err != nil {
		slog.Error(fmt.Sprintf("error reading config: %v", err))
		return
	}

	holder, err := impl.NewHolder(config.DatabaseUrl)
	if err != nil {
		slog.Error(fmt.Sprintf("error creating holder: %v", err))
		return
	}

	app := fiber.New()

	prometheus := fiberprometheus.New("cubepanion-api")
	prometheus.RegisterAt(app, "/metrics")

	app.Use(prometheus.Middleware)
	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{
		Max:               10,
		Expiration:        10 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))
	app.Use(cache.New(cache.Config{
		Storage: redisCache(config.RedisConfig),
		Next: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.Path(), "/ws/")
		},
	}))
	app.Use(impl.SetHolderMiddelware(holder))

	routes.ChestApi(app)
	routes.MapApi(app)
	routes.LeaderboardApi(app)

	// Web Socket Integration
	app.Use("/ws", integration.RequireUpgrade)
	app.Get("/ws/:uuid", integration.Handler())

	err = app.Listen(fmt.Sprintf("%s:%d", config.Address, config.Port))
	if err != nil {
		slog.Error(fmt.Sprintf("error starting server: %v", err))
	}

}
