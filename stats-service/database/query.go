package database

import "database/sql"

var (
	getGameStat *sql.Stmt
	getAllStats *sql.Stmt
)

func load(db *sql.DB) error {
	var err error

	getGameStat, err = db.Prepare(`
		SELECT game,player_count,unix_time_stamp
			FROM game_stats
			WHERE unix_time_stamp =
			      (SELECT MAX(unix_time_stamp)
			       	FROM game_stats
			       	WHERE game = $1)
			  AND game = $1`)
	if err != nil {
		return err
	}

	getAllStats, err = db.Prepare(`
	SELECT game,player_count,unix_time_stamp
	    		FROM game_stats
	    		WHERE unix_time_stamp = ANY(
	    		    SELECT MAX(unix_time_stamp)
	    		    FROM game_stats
	    		    GROUP BY game
	    		))`)
	if err != nil {
		return err
	}

	return nil
}
