package main

import (
	"github.com/Fesaa/CubepanionAPI/core/errors"
	"github.com/Fesaa/CubepanionAPI/core/log"
	"net/url"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/games-service/database"
	"github.com/gofiber/fiber/v2"
)

func games(ms core.MicroService[core.MicroServiceConfig, database.Database], c *fiber.Ctx) error {
	activeS := c.Params("active", "true")
	active := activeS == "true"

	games, err := ms.DB().GetGames(active)
	if err != nil {
		log.Error("Error getting games: ", "error", err)
		return c.Status(500).JSON(errors.AsFiberMap(errors.DBError))
	}

	return c.JSON(games)
}

func game(ms core.MicroService[core.MicroServiceConfig, database.Database], c *fiber.Ctx) error {
	game := c.Params("game")
	if game == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "game name is required",
		})
	}

	var err error
	game, err = url.QueryUnescape(game)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "game name is invalid",
		})
	}

	g, err := ms.DB().GetGame(game)
	if err != nil {
		log.Error("Error getting game: ", "error", err)
		return c.Status(500).JSON(errors.AsFiberMap(errors.DBError))
	}

	return c.SendString(g)
}
