package database

import (
	"database/sql"
	"fmt"
)

var getLastSubmission *sql.Stmt
var getGameLastSubmissions *sql.Stmt
var getLeaderboard *sql.Stmt

var newSubmission *sql.Stmt
var disableSubmission *sql.Stmt

func load(db *sql.DB) error {
	var err error

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
