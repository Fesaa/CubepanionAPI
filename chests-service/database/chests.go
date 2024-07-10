package database

import (
	"database/sql"
	"github.com/Fesaa/CubepanionAPI/core/log"
	"github.com/Fesaa/CubepanionAPI/core/models"
)

func innerGetChestLocations(season string) ([]models.ChestLocation, error) {
	rows, err := getChestLocations.Query(season)
	if err != nil {
		log.Error("Error querying for chest locations: ", "error", err)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Warn("Error closing rows: ", "error", err)
		}
	}(rows)
	var chestLocations []models.ChestLocation
	for rows.Next() {
		var cl models.ChestLocation
		err = rows.Scan(&cl.Season, &cl.X, &cl.Y, &cl.Z)
		if err != nil {
			log.Error("Error scanning chest location:", "error", err)
			return nil, err
		}
		chestLocations = append(chestLocations, cl)
	}

	return chestLocations, nil
}

func innerGetCurrentChestLocations() ([]models.ChestLocation, error) {
	rows, err := getCurrentChestLocations.Query()
	if err != nil {
		log.Error("Error querying for current chest locations: ", "error", err)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Warn("Error closing rows: ", "error", err)
		}
	}(rows)
	var chestLocations []models.ChestLocation
	for rows.Next() {
		var cl models.ChestLocation
		err = rows.Scan(&cl.Season, &cl.X, &cl.Y, &cl.Z)
		if err != nil {
			log.Error("Error scanning chest location: ", "error", err)
			return nil, err
		}
		chestLocations = append(chestLocations, cl)
	}

	return chestLocations, nil
}

func innerGetSeasons(active bool) ([]models.Season, error) {
	var rows *sql.Rows
	var err error

	if active {
		rows, err = getActiveSeasons.Query()
	} else {
		rows, err = getSeasons.Query()
	}

	if err != nil {
		log.Error("Error querying for seasons: ", "error", err)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Warn("Error closing rows: ", "error", err)
		}
	}(rows)
	var seasons []models.Season
	for rows.Next() {
		var season models.Season
		err = rows.Scan(&season)
		if err != nil {
			log.Error("Error scanning season: ", "error", err)
			return nil, err
		}
		seasons = append(seasons, season)
	}

	return seasons, nil
}
