package database

import (
	"database/sql"
	"github.com/Fesaa/CubepanionAPI/core/log"
	"github.com/Fesaa/CubepanionAPI/core/models"
)

func innerGetEggWarsMaps() ([]models.GameMap, error) {
	maps, err := getGameMaps.Query()
	if err != nil {
		log.Error("Error querying for game maps", "errors", err)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Error("Error closing rows: ", "error", err)
		}
	}(maps)

	var ret []models.GameMap
	for maps.Next() {
		var gm models.GameMap
		err = maps.Scan(&gm.Game, &gm.UniqueName, &gm.MapName, &gm.TeamSize, &gm.BuildLimit, &gm.Colours, &gm.Layout)
		if err != nil {
			log.Error("Error scanning game map", "errors", err)
			return nil, err
		}

		ret = append(ret, gm)
	}

	return ret, nil
}
