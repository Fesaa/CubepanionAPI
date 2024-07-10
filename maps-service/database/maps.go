package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Fesaa/CubepanionAPI/core/log"
	"github.com/Fesaa/CubepanionAPI/core/models"
)

func innerGetEggWarsMap(name string) (*models.EggWarsMap, error) {
	em := getEggWarsMap.QueryRow(name)

	var emap models.EggWarsMap
	err := em.Scan(&emap.UniqueName, &emap.MapName, &emap.TeamSize, &emap.BuildLimit, &emap.Colours, &emap.Layout)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("No eggwars map with unique name %s", name)
		}
		log.Error("Error scanning eggwars map", "error", err)
		return nil, err
	}

	emap.Generators = make([]models.Generator, 0)
	gens, err := getMapGenerators.Query(name)
	if err != nil {
		log.Error("Error querying for generators", "error", err)
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
			log.Error("Error scanning generator", "error", err)
			return nil, err
		}
		emap.Generators = append(emap.Generators, eg)
	}

	return &emap, nil
}

func innerGetEggWarsMaps() ([]models.EggWarsMap, error) {
	maps, err := getEggWarsMaps.Query()
	if err != nil {
		log.Error("Error querying for eggwars maps", "errors", err)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Error("Error closing rows: ", "error", err)
		}
	}(maps)
	var mapsMap = make(map[string]models.EggWarsMap)

	for maps.Next() {
		var em models.EggWarsMap
		err = maps.Scan(&em.UniqueName, &em.MapName, &em.TeamSize, &em.BuildLimit, &em.Colours, &em.Layout)
		if err != nil {
			log.Error("Error scanning eggwars map", "errors", err)
			return nil, err
		}

		em.Generators = make([]models.Generator, 0)
		mapsMap[em.UniqueName] = em
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

	ret := make([]models.EggWarsMap, 0, len(mapsMap))
	for _, v := range mapsMap {
		ret = append(ret, v)
	}

	return ret, nil
}
