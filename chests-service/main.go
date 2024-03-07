package main

import (
	"log/slog"

	"github.com/Fesaa/CubepanionAPI/chests-service/database"
	"github.com/Fesaa/CubepanionAPI/core"
)

func main() {
	config, err := core.LoadDefaultConfig("config.yaml")
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return
	}

	ms, err := core.NewMicroService(config, database.Connect)
	if err != nil {
		slog.Error("Failed to create microservice: ", "error", err)
		return
	}

	ms.UseDefaults()
	ms.UseRedisCache()

	ms.Get("/chests", CurrentChestLocations)
	ms.Get("/chests/:season", ChestLocations)
	ms.Get("/seasons/:active", Seasons)

	err = ms.Start()
	if err != nil {
		slog.Error("Failed to start microservice: ", "error", err)
	}
}
