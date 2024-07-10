package main

import (
	"github.com/Fesaa/CubepanionAPI/core/log"
	"io"
	"net/http"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/leaderboard-service/database"
)

func getGame(ms core.MicroService[LeaderboardServiceConfig, database.Database], s string) (string, error) {
	url := ms.Config().GamesService()
	url += "/game/" + s

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Warn("Error closing response body", "error", err)
		}
	}(res.Body)
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
