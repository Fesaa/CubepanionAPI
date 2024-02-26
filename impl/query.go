package impl

import (
	"database/sql"
	"fmt"
)

var getSeasons *sql.Stmt
var getActiveSeasons *sql.Stmt
var getCurrentChestLocations *sql.Stmt
var getChestLocations *sql.Stmt

var getEggWarsMaps *sql.Stmt
var getGenerators *sql.Stmt
var getEggWarsMap *sql.Stmt
var getMapGenerators *sql.Stmt

var getGames *sql.Stmt
var getActiveGames *sql.Stmt

var getGame *sql.Stmt
var getLastSubmission *sql.Stmt
var getGameLastSubmissions *sql.Stmt
var getLeaderboard *sql.Stmt
var getLeaderboardForPlayer *sql.Stmt

var newSubmission *sql.Stmt
var disableSubmission *sql.Stmt

func loadQueries(db *sql.DB) error {
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

	getGames, err = db.Prepare("SELECT game,active,display_name,aliases,score_type FROM games")
	if err != nil {
		return fmt.Errorf("Error preparing getGames: %v", err)
	}

	getActiveGames, err = db.Prepare("SELECT game,active,display_name,aliases,score_type FROM games WHERE active = true")
	if err != nil {
		return fmt.Errorf("Error preparing getActiveGames: %v", err)
	}

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

	getLastSubmission, err = db.Prepare("SELECT MAX(unix_time_stamp) FROM submissions WHERE game = $1 AND valid = true")
	if err != nil {
		return fmt.Errorf("Error preparing lastSubmission: %v", err)
	}

	getGameLastSubmissions, err = db.Prepare("SELECT game, MAX(unix_time_stamp) FROM submissions WHERE valid = true GROUP BY game")
	if err != nil {
		return fmt.Errorf("Error preparing getGameLastSubmissions: %v", err)
	}

	getLeaderboard, err = db.Prepare(`
		SELECT *
		FROM leaderboards
		WHERE unix_time_stamp = $1
			AND game = $2
			AND position >= $3
			AND position <= $4
		ORDER BY position`)
	if err != nil {
		return fmt.Errorf("Error preparing getLeaderboard: %v", err)
	}

	getLeaderboardForPlayer, err = db.Prepare(`
		SELECT *
		FROM leaderboards
		WHERE UPPER(player) = UPPER($1)
			AND (game, unix_time_stamp) = ANY($2)
		ORDER BY position`)
	if err != nil {
		return fmt.Errorf("Error preparing getLeaderboardForPlayer: %v", err)
	}

	newSubmission, err = db.Prepare("INSERT INTO submissions (uuid, unix_time_stamp, game, valid) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("Error preparing newSubmission: %v", err)
	}

	disableSubmission, err = db.Prepare("UPDATE submissions SET valid = false WHERE uuid = $1 AND unix_time_stamp = $2")
	if err != nil {
		return fmt.Errorf("Error preparing disableSubmission: %v", err)
	}
	return nil
}
