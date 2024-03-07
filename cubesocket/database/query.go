package database

import "database/sql"

var (
	getPlayerLocation    *sql.Stmt
	setPlayerLocation    *sql.Stmt
	removePlayerLocation *sql.Stmt
	getSharedPlayers     *sql.Stmt
)

func load(db *sql.DB) error {
	var err error

	getPlayerLocation, err = db.Prepare("SELECT current, previous, in_pre_lobby FROM player_locations WHERE uuid = $1")
	if err != nil {
		return err
	}

	setPlayerLocation, err = db.Prepare("INSERT INTO player_locations (uuid, current, previous, in_pre_lobby) VALUES ($1, $2, $3, $4) ON CONFLICT (uuid) DO UPDATE SET current = $2, previous = $3, in_pre_lobby = $4")
	if err != nil {
		return err
	}

	removePlayerLocation, err = db.Prepare("DELETE FROM player_locations WHERE uuid = $1")
	if err != nil {
		return err
	}

	getSharedPlayers, err = db.Prepare("SELECT uuid FROM player_locations WHERE current = (SELECT current FROM player_locations WHERE uuid = $1) AND uuid != $1")
	if err != nil {
		return err
	}

	return nil
}
