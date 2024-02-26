package impl

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Fesaa/CubepanionAPI/models"
)

func innerGetEggWarsMap(name string) (*models.EggWarsMap, error) {
	em := getEggWarsMap.QueryRow(name)

	var emap models.EggWarsMap
	err := em.Scan(&emap.UniqueName, &emap.MapName, &emap.TeamSize, &emap.BuildLimit, &emap.Colours, &emap.Layout)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No eggwars map with unique name %s", name)
		}
		slog.Error("Error scanning eggwars map: %v", err)
		return nil, err
	}

	emap.Generators = make([]models.Generator, 0)
	gens, err := getMapGenerators.Query(name)
	if err != nil {
		slog.Error("Error querying for generators: %v", err)
		return nil, err
	}

	defer gens.Close()
	for gens.Next() {
		var eg models.Generator
		err = gens.Scan(&eg.UniqueName, &eg.Ordering, &eg.Type, &eg.Location, &eg.Level, &eg.Count)
		if err != nil {
			slog.Error("Error scanning generator: %v", err)
			return nil, err
		}
		emap.Generators = append(emap.Generators, eg)
	}

	return &emap, nil
}

func innerGetEggWarsMaps() ([]models.EggWarsMap, error) {
	maps, err := getEggWarsMaps.Query()
	if err != nil {
		slog.Error("Error querying for eggwars maps: %v", err)
		return nil, err
	}

	defer maps.Close()
	var mapsMap map[string]models.EggWarsMap = make(map[string]models.EggWarsMap)

	for maps.Next() {
		var em models.EggWarsMap
		err = maps.Scan(&em.UniqueName, &em.MapName, &em.TeamSize, &em.BuildLimit, &em.Colours, &em.Layout)
		if err != nil {
			slog.Error("Error scanning eggwars map: %v", err)
			return nil, err
		}

		em.Generators = make([]models.Generator, 0)
		mapsMap[em.UniqueName] = em
	}

	gens, err := getGenerators.Query()
	if err != nil {
		slog.Error("Error querying for generators: %v", err)
		return nil, err
	}

	defer gens.Close()
	for gens.Next() {
		var eg models.Generator
		err = gens.Scan(&eg.UniqueName, &eg.Ordering, &eg.Type, &eg.Location, &eg.Level, &eg.Count)
		if err != nil {
			slog.Error("Error scanning generator: %v", err)
			return nil, err
		}

		em, ok := mapsMap[eg.UniqueName]
		if !ok {
			slog.Error(fmt.Sprintf("Generator %s does not have a corresponding map", eg.UniqueName))
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
