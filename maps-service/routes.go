package main

import (
	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/core/errors"
	"github.com/Fesaa/CubepanionAPI/core/log"
	"github.com/Fesaa/CubepanionAPI/maps-service/database"
	"github.com/gofiber/fiber/v2"
)

func Maps(ms core.MicroService[core.MicroServiceConfig, database.Database], c *fiber.Ctx) error {
	maps, err := ms.DB().GetAllMaps()
	if err != nil {
		log.Error("Error getting maps", "error", err)
		return c.Status(500).JSON(errors.AsFiberMap(errors.DBError))
	}

	return c.JSON(maps)
}

func Map(ms core.MicroService[core.MicroServiceConfig, database.Database], c *fiber.Ctx) error {
	name := c.Params("mapName")
	maps, err := ms.DB().GetMap(name)
	if err != nil {
		log.Error("Error getting map", "error", err)
		return c.Status(500).JSON(errors.AsFiberMap(errors.DBError))
	}

	return c.JSON(maps)
}
