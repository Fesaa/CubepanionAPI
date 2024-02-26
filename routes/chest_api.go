package routes

import (
	"github.com/Fesaa/CubepanionAPI/models"
	"github.com/gofiber/fiber/v2"
)

func ChestApi(app *fiber.App) {
	g := app.Group("/chest_api")
	g.Get("/current", chestAPI_current)
	g.Get("/seasons/:running", chestAPI_seasons)
	g.Get("/season/:season", chestAPI_season)
}

func chestAPI_current(c *fiber.Ctx) error {
	holder, _ := c.Locals(models.HOLDER_KEY).(models.Holder)
	db := holder.GetDatabaseProvider()
	locations, err := db.GetCurrentChestLocations()
	if err != nil {
		return err
	}
	return c.JSON(locations)
}

func chestAPI_seasons(c *fiber.Ctx) error {
	holder, _ := c.Locals(models.HOLDER_KEY).(models.Holder)
	db := holder.GetDatabaseProvider()

	runningS := c.Params("running", "false")
	running := runningS == "true"

	seasons, err := db.GetSeasons(running)
	if err != nil {
		return err
	}
	return c.JSON(seasons)
}

func chestAPI_season(c *fiber.Ctx) error {
	holder, _ := c.Locals(models.HOLDER_KEY).(models.Holder)
	db := holder.GetDatabaseProvider()

	season := c.Params("season")
	seasons, err := db.GetChestLocations(season)
	if err != nil {
		return err
	}
	return c.JSON(seasons)
}
