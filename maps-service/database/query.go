package database

import (
	"database/sql"
	"fmt"
)

var getGameMaps *sql.Stmt

func load(db *sql.DB) error {
	var err error

	getGameMaps, err = db.Prepare("SELECT game,unique_name,map_name,team_size,build_limit,colours,layout  FROM game_maps")
	if err != nil {
		return fmt.Errorf("Error preparing getGameMaps: %v", err)
	}

	return nil
}
