package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

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
func generateLeaderboardInsertSQL(unix uint64, game string, rows []models.LeaderboardRow) string {
	s := "INSERT INTO leaderboards (game, player, normalized_player_name, position, score, unix_time_stamp) VALUES "
	for i, row := range rows {
		if i != 0 {
			s += ", "
		}
		s += fmt.Sprintf("('%s', '%s', '%s', %d, %d, %d)", game, row.Player, strings.ToUpper(row.Player), row.Position, row.Score, unix)
	}

	return s
}

func innerInsertLeaderboards(db *sql.DB, req models.LeaderboardSubmission) error {
	if len(req.Entries) != models.LEADERBOARD_SIZE {
		return fmt.Errorf("Leaderboard submission must have 200 entries")
	}

	err := innerInsertSubmission(req.Uuid, req.Game, req.UnixTimeStamp)
	if err != nil {
		return err
	}

	_, err = db.Exec(generateLeaderboardInsertSQL(req.UnixTimeStamp, req.Game, req.Entries))
	if err != nil {
		err2 := innerDisableSubmission(req.Uuid, req.UnixTimeStamp)
		if err2 != nil {
			slog.Error(fmt.Sprintf("Error disabling submission after failed leaderboard insert: %v", err2))
		}
		return err
	}

	return nil
}
