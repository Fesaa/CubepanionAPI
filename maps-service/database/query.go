package database

import (
	"database/sql"
	"fmt"
)

var getEggWarsMaps *sql.Stmt
var getGenerators *sql.Stmt
var getEggWarsMap *sql.Stmt
var getMapGenerators *sql.Stmt

func load(db *sql.DB) error {
	var err error

	getEggWarsMaps, err = db.Prepare("SELECT * FROM eggwars_maps")
	if err != nil {
		return fmt.Errorf("Error preparing getEggWarsMaps: %v", err)
	}

	getGenerators, err = db.Prepare("SELECT * FROM generators")
	if err != nil {
		return fmt.Errorf("Error preparing getGenerators: %v", err)
	}

	getEggWarsMap, err = db.Prepare("SELECT * FROM eggwars_maps WHERE unique_name = $1")
	if err != nil {
		return fmt.Errorf("Error preparing getEggWarsMap: %v", err)
	}

	getMapGenerators, err = db.Prepare("SELECT * FROM generators WHERE unique_name = $1")
	if err != nil {
		return fmt.Errorf("Error preparing getMapGenerators: %v", err)
	}


	return nil
}
