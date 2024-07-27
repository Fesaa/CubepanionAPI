package database

import (
	"database/sql"
	"fmt"
	"github.com/Fesaa/CubepanionAPI/core/log"
	"github.com/Fesaa/CubepanionAPI/core/models"
)

func innerGetEggWarsMaps() ([]models.GameMap, error) {
	maps, err := getGameMaps.Query()
	if err != nil {
		log.Error("Error querying for eggwars maps", "errors", err)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Error("Error closing rows: ", "error", err)
		}
	}(maps)
	var mapsMap = make(map[string]models.GameMap)

	for maps.Next() {
		var gm models.GameMap
		err = maps.Scan(&gm.Game, &gm.UniqueName, &gm.MapName, &gm.TeamSize, &gm.BuildLimit, &gm.Colours, &gm.Layout)
		if err != nil {
			log.Error("Error scanning eggwars map", "errors", err)
			return nil, err
		}

		mapsMap[gm.UniqueName] = gm
	}

	gens, err := getGenerators.Query()
	if err != nil {
		log.Error("Error querying for generators", "errors", err)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Error("Error closing rows: ", "error", err)
		}
	}(gens)
	for gens.Next() {
		var eg models.Generator
		err = gens.Scan(&eg.UniqueName, &eg.Ordering, &eg.Type, &eg.Location, &eg.Level, &eg.Count)
		if err != nil {
			log.Error("Error scanning generator", "errors", err)
			return nil, err
		}

		em, ok := mapsMap[eg.UniqueName]
		if !ok {
			log.Error(fmt.Sprintf("Generator %s does not have a corresponding map", eg.UniqueName))
			continue
		}
		em.Generators = append(em.Generators, eg)
		mapsMap[eg.UniqueName] = em
	}

	ret := make([]models.GameMap, 0, len(mapsMap))
	for _, v := range mapsMap {
		ret = append(ret, v)
	}

	return ret, nil
}
