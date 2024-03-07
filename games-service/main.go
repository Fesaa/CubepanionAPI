package main

import (
	"log/slog"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/games-service/database"
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

	ms.Get("/games/:active", games)
	ms.Get("/game/:game", game)

	err = ms.Start()
	if err != nil {
		slog.Error("Failed to start microservice: ", "error", err)
	}
}
