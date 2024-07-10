package main

import (
	"github.com/Fesaa/CubepanionAPI/chests-service/database"
	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/core/errors"
	"github.com/Fesaa/CubepanionAPI/core/log"
	"github.com/gofiber/fiber/v2"
)

func Seasons(ms core.MicroService[core.MicroServiceConfig, database.Database], c *fiber.Ctx) error {
	activeS := c.Params("active", "false")
	active := activeS == "true"
	seasons, err := ms.DB().GetSeasons(active)
	if err != nil {
		log.Error("Error getting seasons: ", "error", err)
		return c.Status(500).JSON(errors.AsFiberMap(errors.DBError))
	}
	return c.JSON(seasons)
}

func ChestLocations(ms core.MicroService[core.MicroServiceConfig, database.Database], c *fiber.Ctx) error {
	season := c.Params("season")
	seasons, err := ms.DB().GetChests(season)
	if err != nil {
		log.Error("Error getting chest locations: ", "error", err)
		return c.Status(500).JSON(errors.AsFiberMap(errors.DBError))
	}
	return c.JSON(seasons)
}

func CurrentChestLocations(ms core.MicroService[core.MicroServiceConfig, database.Database], c *fiber.Ctx) error {
	seasons, err := ms.DB().GetCurrentChests()
	if err != nil {
		log.Error("Error getting current chest locations: ", "error", err)
		return c.Status(500).JSON(errors.AsFiberMap(errors.DBError))
	}
	return c.JSON(seasons)
}
