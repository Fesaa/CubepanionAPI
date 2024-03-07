package database

import (
	"database/sql"
	"fmt"
)

var getSeasons *sql.Stmt
var getActiveSeasons *sql.Stmt
var getCurrentChestLocations *sql.Stmt
var getChestLocations *sql.Stmt

func load(db *sql.DB) error {
	var err error

	getSeasons, err = db.Prepare("SELECT season_name FROM seasons")
	if err != nil {
		return fmt.Errorf("Error preparing getSeasons: %v", err)
	}

	getActiveSeasons, err = db.Prepare("SELECT season_name FROM seasons WHERE running = true")
	if err != nil {
		return fmt.Errorf("Error preparing getActiveSeasons: %v", err)
	}

	getCurrentChestLocations, err = db.Prepare("SELECT season_name, x, y, z FROM chest_locations WHERE season_name IN (SELECT season_name FROM seasons WHERE running = true)")
	if err != nil {
		return fmt.Errorf("Error preparing getCurrentChestLocations: %v", err)
	}

	getChestLocations, err = db.Prepare("SELECT season_name, x, y, z FROM chest_locations WHERE season_name = $1")
	if err != nil {
		return fmt.Errorf("Error preparing getChestLocations: %v", err)
	}

	return nil
}
