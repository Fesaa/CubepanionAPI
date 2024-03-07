package main

import (
	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/maps-service/database"
	"github.com/gofiber/fiber/v2"
)

func Maps(ms core.MicroService[core.MicroServiceConfig], c *fiber.Ctx) error {
	maps, err := database.GetAllMaps()
	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(maps)
}

func Map(ms core.MicroService[core.MicroServiceConfig], c *fiber.Ctx) error {
	name := c.Params("mapName")
	maps, err := database.GetMap(name)
	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(maps)
}
