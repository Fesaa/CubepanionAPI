package database

import (
	"database/sql"
	"fmt"
)

var getLeaderboard *sql.Stmt
var getLeaderboardForPlayer *sql.Stmt
var getLeaderboardForPlayers *sql.Stmt

var newSubmission *sql.Stmt
var disableSubmission *sql.Stmt

func load(db *sql.DB) error {
	var err error

	getLeaderboard, err = db.Prepare(`
		SELECT game, player, position, score, unix_time_stamp
		FROM leaderboards
		WHERE unix_time_stamp = (
			SELECT MAX(unix_time_stamp)
			FROM submissions
			WHERE valid = true
				AND game = $1
			)
			AND game = $1
			AND position >= $2
			AND position <= $3
		ORDER BY position`)
	if err != nil {
		return fmt.Errorf("Error preparing getLeaderboard: %v", err)
	}

	getLeaderboardForPlayer, err = db.Prepare(`
		SELECT game, player, position, score, unix_time_stamp
		FROM leaderboards
		WHERE normalized_player_name = UPPER($1)
		AND (unix_time_stamp, game) = ANY(
			SELECT MAX(unix_time_stamp), game
			FROM submissions
			WHERE valid = true
			GROUP BY game
			)
		ORDER BY position
		`)
	if err != nil {
		return fmt.Errorf("Error preparing getLeaderboardForPlayer: %v", err)
	}

	getLeaderboardForPlayers, err = db.Prepare(`
		SELECT game, player, position, score, unix_time_stamp
		FROM leaderboards
		WHERE game = $1
			AND normalized_player_name = ANY(SELECT UPPER(unnest($2::text[])))
			AND unix_time_stamp = (
				SELECT MAX(unix_time_stamp)
				FROM submissions
				WHERE game = $1
				AND valid = true
				)
		ORDER BY position`)
	if err != nil {
		return fmt.Errorf("Error preparing getLeaderboardForPlayers: %v", err)
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
