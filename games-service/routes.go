package main

import (
	"net/url"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/games-service/database"
	"github.com/gofiber/fiber/v2"
)

func games(ms core.MicroService[core.MicroServiceConfig], c *fiber.Ctx) error {
	activeS := c.Params("active", "true")
	active := activeS == "true"

	games, err := database.GetGames(active)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(games)
}

func game(ms core.MicroService[core.MicroServiceConfig], c *fiber.Ctx) error {
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

	g, err := database.GetGame(game)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendString(g)
}
