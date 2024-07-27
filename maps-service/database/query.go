package database

import (
	"database/sql"
	"fmt"
)

var getGameMaps *sql.Stmt
var getGenerators *sql.Stmt

func load(db *sql.DB) error {
	var err error

	getGameMaps, err = db.Prepare("SELECT * FROM game_maps")
	if err != nil {
		return fmt.Errorf("Error preparing getGameMaps: %v", err)
	}

	getGenerators, err = db.Prepare("SELECT * FROM generators")
	if err != nil {
		return fmt.Errorf("Error preparing getGenerators: %v", err)
	}

	return nil
}
