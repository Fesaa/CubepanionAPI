package main

import (
	"log/slog"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/leaderboard-service/database"
)

func main() {
	var config LeaderboardServiceConfig = LeaderboardServiceConfig{}
	err := core.LoadConfig("config.yaml", &config)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return
	}

	ms, err := core.NewMicroService(config, database.Connect)

	ms.UseDefaults()
	ms.UseRedisCache()

	ms.Post("/", Submit)
	ms.Get("/player/:name", PlayerLeaderboard)
	ms.Get("/game/:game", GameLeaderboard)
	ms.Get("/game/:game/bounded", GameLeaderboardBounded)

	err = ms.Start()
	if err != nil {
		slog.Error("Failed to start microservice", "error", err)
	}
}
