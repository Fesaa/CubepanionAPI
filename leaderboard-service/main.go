package main

import (
	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/core/log"
	"github.com/Fesaa/CubepanionAPI/leaderboard-service/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func main() {
	var config LeaderboardServiceConfig = LeaderboardServiceConfig{}
	err := core.LoadConfig("config.yaml", &config)
	if err != nil {
		log.Error("Failed to load config", "error", err)
		return
	}

	ms, err := core.NewMicroService(config, database.Connect)
	if err != nil {
		log.Error("Failed to create microservice: ", "error", err)
		return
	}

	ms.UseDefaults()
	ms.UseRedisCache(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Method() != "GET" && c.Path() != "/batch"
		},
		Methods: []string{fiber.MethodGet, fiber.MethodPost},
	})

	ms.Post("/", Submit)
	ms.Get("/player/:name", PlayerLeaderboard)
	ms.Get("/game/:game", GameLeaderboard)
	ms.Get("/game/:game/bounded", GameLeaderboardBounded)
	ms.Post("/batch", BatchPlayerLeaderboard)

	err = ms.Start()
	if err != nil {
		log.Error("Failed to start microservice", "error", err)
	}
}
