package database

import (
	"database/sql"
	"fmt"
	"github.com/Fesaa/CubepanionAPI/core/log"
	"strings"
	"time"

	"github.com/Fesaa/CubepanionAPI/core/models"
)

func innerInsertSubmission(uuid, game string, unix uint64) error {
	_, err := newSubmission.Exec(uuid, unix, game, true)
	return err
}

func innerDisableSubmission(uuid string, unix uint64) error {
	_, err := disableSubmission.Exec(uuid, unix)
	return err
}

// Aware this is rather unsafe, but we can assume it's fine as game is provided by the server
// And player has been checked against a regex
func generateLeaderboardInsertSQL(unix int64, game string, rows []models.LeaderboardRow) string {
	s := "INSERT INTO leaderboards (game, player, normalized_player_name, position, score, texture, unix_time_stamp) VALUES "
	for i, row := range rows {
		if i != 0 {
			s += ", "
		}

		var texture string
		if row.Texture == "" {
			texture = "NULL"
		} else {
			texture = fmt.Sprintf("'%s'", row.Texture)
		}

		s += fmt.Sprintf("('%s', '%s', '%s', %d, %d, %s, %d)", game, row.Player, strings.ToUpper(row.Player), row.Position, row.Score, texture, unix)
	}

	return s
}

func innerInsertLeaderboards(db *sql.DB, req models.LeaderboardSubmission) error {
	if len(req.Entries) != models.LEADERBOARD_SIZE {
		return fmt.Errorf("leaderboard submission must have 200 entries")
	}

	err := innerInsertSubmission(req.Uuid, req.Game, req.UnixTimeStamp)
	if err != nil {
		return err
	}

	_, err = db.Exec(generateLeaderboardInsertSQL(time.Now().UnixMilli(), req.Game, req.Entries))
	if err != nil {
		err2 := innerDisableSubmission(req.Uuid, req.UnixTimeStamp)
		if err2 != nil {
			log.Error("Error disabling submission after failed leaderboard insert", "error", err2)
		}
		return err
	}

	return nil
}
