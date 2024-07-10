package main

import (
	"log/slog"
	"time"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/maps-service/database"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func main() {
	config, err := core.LoadDefaultConfig("config.yaml")
	if err != nil {
		slog.Error("error reading config: ", "error", err)
		return
	}

	ms, err := core.NewMicroService(config, database.Connect)
	if err != nil {
		slog.Error("error creating microservice: ", "error", err)
		return
	}

	ms.UseDefaults()
	ms.UseRedisCache(cache.Config{
		Expiration: 1 * time.Hour,
	})

	ms.Get("/", Maps)
	ms.Get("/:mapName", Map)

	err = ms.Start()
	if err != nil {
		slog.Error("error starting microservice: ", "error", err)
	}
}
