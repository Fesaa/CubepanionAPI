package main

import (
	"github.com/Fesaa/CubepanionAPI/chests-service/database"
	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/gofiber/fiber/v2"
)

func Seasons(ms core.MicroService[core.MicroServiceConfig, database.Database], c *fiber.Ctx) error {
	activeS := c.Params("active", "false")
	active := activeS == "true"
	seasons, err := ms.DB().GetSeasons(active)
	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(seasons)
}

func ChestLocations(ms core.MicroService[core.MicroServiceConfig, database.Database], c *fiber.Ctx) error {
	season := c.Params("season")
	seasons, err := ms.DB().GetChests(season)
	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(seasons)
}

func CurrentChestLocations(ms core.MicroService[core.MicroServiceConfig, database.Database], c *fiber.Ctx) error {
	seasons, err := ms.DB().GetCurrentChests()
	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(seasons)
}
