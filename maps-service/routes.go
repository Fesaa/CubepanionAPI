package main

import (
	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/maps-service/database"
	"github.com/gofiber/fiber/v2"
)

func Maps(ms core.MicroService[core.MicroServiceConfig, database.Database], c *fiber.Ctx) error {
	maps, err := ms.DB().GetAllMaps()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(maps)
}

func Map(ms core.MicroService[core.MicroServiceConfig, database.Database], c *fiber.Ctx) error {
	name := c.Params("mapName")
	maps, err := ms.DB().GetMap(name)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(maps)
}
