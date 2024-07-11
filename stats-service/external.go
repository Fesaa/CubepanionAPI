package main

import (
	"fmt"
	"github.com/Fesaa/CubepanionAPI/core/log"
	"io"
	"net/http"
	"regexp"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/stats-service/database"
)

var gameRegex = regexp.MustCompile(`^[a-zA-Z0-9_ ]`)

func convertGame(ms core.MicroService[StatsServiceConfig, database.Database], game string) (string, error) {
	if game == "" {
		return "", fmt.Errorf("game name is required")
	}

	if !gameRegex.MatchString(game) {
		return "", fmt.Errorf("game name is invalid")
	}

	gameDisplayName, err := getGame(ms, game)
	if err != nil {
		return "", fmt.Errorf("game %s does not exist", game)
	}

	return gameDisplayName, nil
}

func getGame(ms core.MicroService[StatsServiceConfig, database.Database], s string) (string, error) {
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
