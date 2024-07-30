package main

import (
	"github.com/Fesaa/CubepanionAPI/core/log"
	"time"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/games-service/database"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func main() {
	config, err := core.LoadDefaultConfig("config.yaml")
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
		Expiration:   1 * time.Hour,
		CacheControl: true,
	})

	ms.Get("/:active", games)
	ms.Get("/game/:game", game)

	err = ms.Start()
	if err != nil {
		log.Error("Failed to start microservice: ", "error", err)
	}
}
