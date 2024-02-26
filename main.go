package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Fesaa/CubepanionAPI/impl"
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
		Expiration:        60 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))
	app.Use(cache.New(cache.Config{
		Storage: redisCache(config.RedisConfig),
	}))
	app.Use(impl.SetHolderMiddelware(holder))

	routes.ChestApi(app)
	routes.MapApi(app)
	routes.LeaderboardApi(app)

	err = app.Listen(fmt.Sprintf("%s:%d", config.Address, config.Port))
	if err != nil {
		slog.Error(fmt.Sprintf("error starting server: %v", err))
	}

}
