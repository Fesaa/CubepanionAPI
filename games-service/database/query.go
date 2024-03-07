package database

import (
	"database/sql"
	"fmt"
)

var getGame *sql.Stmt
var getGames *sql.Stmt
var getActiveGames *sql.Stmt

func load(db *sql.DB) error {
	var err error

	getGame, err = db.Prepare(`
		SELECT display_name
		FROM games
		WHERE UPPER(game) = UPPER($1)
			OR UPPER(display_name) = UPPER($1)
			OR $1 IN (SELECT unnest(string_to_array(aliases, ',')))
		`)
	if err != nil {
		return fmt.Errorf("Error preparing getGame: %v", err)
	}

	getGames, err = db.Prepare("SELECT game,active,display_name,aliases,score_type FROM games")
	if err != nil {
		return fmt.Errorf("Error preparing getGames: %v", err)
	}

	getActiveGames, err = db.Prepare("SELECT game,active,display_name,aliases,score_type FROM games WHERE active = true")
	if err != nil {
		return fmt.Errorf("Error preparing getActiveGames: %v", err)
	}

	return nil
}
