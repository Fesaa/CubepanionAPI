package routes

import (
	"github.com/Fesaa/CubepanionAPI/models"
	"github.com/gofiber/fiber/v2"
)

func MapApi(app *fiber.App) {
	g := app.Group("/eggwars_map_api")
	g.Get("/", mapAPI_all)
	g.Get("/:mapName", mapAPI_specific)
}

func mapAPI_specific(c *fiber.Ctx) error {
	holder, _ := c.Locals(models.HOLDER_KEY).(models.Holder)
	db := holder.GetDatabaseProvider()

	mapName := c.Params("mapName")
	maps, err := db.GetEggWarsMap(mapName)
	if err != nil {
		return err
	}

	return c.JSON(maps)
}

func mapAPI_all(c *fiber.Ctx) error {
	holder, _ := c.Locals(models.HOLDER_KEY).(models.Holder)
	db := holder.GetDatabaseProvider()

	maps, err := db.GetEggWarsMaps()
	if err != nil {
		return err
	}

	return c.JSON(maps)
}
