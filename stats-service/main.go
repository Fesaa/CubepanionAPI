package main

import (
	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/core/log"
	"github.com/Fesaa/CubepanionAPI/stats-service/database"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func main() {
	var config StatsServiceConfig = StatsServiceConfig{}
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
		CacheControl: true,
	})

	ms.Get("/game/:game", GetGameStat)
	ms.Get("/games", GetAllStats)

	err = ms.Start()
	if err != nil {
		log.Error("Failed to start microservice: ", "error", err)
	}
}
