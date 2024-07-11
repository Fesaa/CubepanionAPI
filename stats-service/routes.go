package main

import (
	"fmt"
	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/core/errors"
	"github.com/Fesaa/CubepanionAPI/core/log"
	"github.com/Fesaa/CubepanionAPI/stats-service/database"
	"github.com/gofiber/fiber/v2"
)

func GetGameStat(ms core.MicroService[StatsServiceConfig, database.Database], ctx *fiber.Ctx) error {
	gameS := ctx.Params("game")
	if gameS == "" {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	game, err := convertGame(ms, gameS)
	if err != nil {
		log.Error("Could not convert game", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Could not convert game %s", gameS),
		})
	}

	stat, err := ms.DB().GetStatsForGame(game)
	if err != nil {
		log.Error("Failed to get game stat", "error", err)
		return ctx.JSON(errors.AsFiberMap(errors.DBError))
	}

	if stat == nil {
		return ctx.JSON(fiber.Map{})
	}

	return ctx.JSON(stat)
}

func GetAllStats(ms core.MicroService[StatsServiceConfig, database.Database], ctx *fiber.Ctx) error {
	stats, err := ms.DB().GetAllStats()
	if err != nil {
		log.Error("Failed to get all stats", "error", err)
		return ctx.JSON(errors.AsFiberMap(errors.DBError))
	}

	return ctx.JSON(stats)
}
