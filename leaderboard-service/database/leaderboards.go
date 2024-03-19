package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Fesaa/CubepanionAPI/core/models"
	"github.com/lib/pq"
)

type gameSubmission struct {
	game           string
	lastSubmission int64
}

func (g gameSubmission) Value() string {
	return fmt.Sprintf("('%s', %d)", g.game, g.lastSubmission)
}

func generatePlayerSQL(player string, lsu []gameSubmission) string {
	s := fmt.Sprintf("SELECT * FROM leaderboards WHERE UPPER(player) = UPPER('%s')", player)
	if len(lsu) > 0 {
		subMap := ""
		for _, sub := range lsu {
			subMap += sub.Value() + ", "
		}
		s += fmt.Sprintf(" AND (game, unix_time_stamp) IN (%s)", strings.TrimSuffix(subMap, ", "))
	}

	return s + "ORDER BY position;"
}

func innerGetLeaderboardForPlayer(db *sql.DB, player string) ([]models.LeaderboardRow, error) {
	rows, err := getGameLastSubmissions.Query()
	if err != nil {
		slog.Error("Error querying for last submissions: %v", err)
		return nil, err
	}

	defer rows.Close()
	var lastSubmissions []gameSubmission = make([]gameSubmission, 0)
	for rows.Next() {
		var game string
		var lastSubmission int64
		err = rows.Scan(&game, &lastSubmission)
		if err != nil {
			slog.Error(fmt.Sprintf("Error scanning last submission: %v", err))
			return nil, err
		}
		lastSubmissions = append(lastSubmissions, gameSubmission{game, lastSubmission})
	}

	var leaderboard []models.LeaderboardRow = make([]models.LeaderboardRow, 0)
	rows, err = db.Query(generatePlayerSQL(player, lastSubmissions))
	if err != nil {
		slog.Error(fmt.Sprintf("Error querying for leaderboard: %v", err))
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var row models.LeaderboardRow
		err = rows.Scan(&row.Game, &row.Player, &row.Position, &row.Score, &row.UnixTimeStamp)
		if err != nil {
			slog.Error(fmt.Sprintf("Error scanning leaderboard row: %v", err))
			return nil, err
		}
		leaderboard = append(leaderboard, row)
	}

	return leaderboard, nil
}

func innerGetLeaderboardBounded(game string, start, end int) ([]models.LeaderboardRow, error) {
	unixRow := getLastSubmission.QueryRow(game)
	var lastSubmission int64
	err := unixRow.Scan(&lastSubmission)
	if err != nil {
		slog.Error("Error scanning last submission: %v", err)
		return nil, err
	}

	rows, err := getLeaderboard.Query(lastSubmission, game, start, end)
	if err != nil {
		slog.Error("Error querying for leaderboard: %v", err)
		return nil, err
	}

	defer rows.Close()
	var leaderboard []models.LeaderboardRow = make([]models.LeaderboardRow, 0)
	for rows.Next() {
		var row models.LeaderboardRow
		err = rows.Scan(&row.Game, &row.Player, &row.Position, &row.Score, &row.UnixTimeStamp)
		if err != nil {
			slog.Error(fmt.Sprintf("Error scanning leaderboard row: %v", err))
			return nil, err
		}
		leaderboard = append(leaderboard, row)
	}
	return leaderboard, nil
}

func innerGetLeaderboardForPlayers(req models.GamePlayersRequest) ([]models.LeaderboardRow, error) {
	rows, err := getLeaderboardForPlayers.Query(req.Game, pq.Array(req.Players))
	if err != nil {
		slog.Error("Error querying for leaderboard for players in a game: ", "error", err)
		return nil, err
	}

	defer rows.Close()
	var leaderboard []models.LeaderboardRow = make([]models.LeaderboardRow, 0)
	for rows.Next() {
		var row models.LeaderboardRow
		err = rows.Scan(&row.Game, &row.Player, &row.Position, &row.Score, &row.UnixTimeStamp)
		if err != nil {
			slog.Error(fmt.Sprintf("Error scanning leaderboard row: ", "error", err))
			return nil, err
		}
		leaderboard = append(leaderboard, row)
	}

	return leaderboard, nil
}
