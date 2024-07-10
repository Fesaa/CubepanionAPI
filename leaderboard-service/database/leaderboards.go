package database

import (
	"database/sql"
	"log/slog"

	"github.com/Fesaa/CubepanionAPI/core/log"
	"github.com/Fesaa/CubepanionAPI/core/models"
	"github.com/lib/pq"
)

func innerGetLeaderboardForPlayer(player string) ([]models.LeaderboardRow, error) {
	var leaderboard = make([]models.LeaderboardRow, 0)
	rows, err := getLeaderboardForPlayer.Query(player)
	if err != nil {
		log.Error("Error querying for leaderboard for player: ", "error, err")
		return nil, err
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Warn("Error closing rows: ", "error", err)
		}
	}(rows)
	for rows.Next() {
		var row models.LeaderboardRow
		err = rows.Scan(&row.Game, &row.Player, &row.Position, &row.Score, &row.UnixTimeStamp)
		if err != nil {
			log.Error("Error scanning leaderboard row", "error", err)
			return nil, err
		}
		leaderboard = append(leaderboard, row)
	}

	return leaderboard, nil
}

func innerGetLeaderboardBounded(game string, start, end int) ([]models.LeaderboardRow, error) {
	rows, err := getLeaderboard.Query(game, start, end)
	if err != nil {
		slog.Error("Error querying for leaderboard: %v", err)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Warn("Error closing rows: ", "error", err)
		}
	}(rows)
	var leaderboard = make([]models.LeaderboardRow, 0)
	for rows.Next() {
		var row models.LeaderboardRow
		err = rows.Scan(&row.Game, &row.Player, &row.Position, &row.Score, &row.UnixTimeStamp)
		if err != nil {
			log.Error("Error scanning leaderboard row", "error", err)
			return nil, err
		}
		leaderboard = append(leaderboard, row)
	}
	return leaderboard, nil
}

func innerGetLeaderboardForPlayers(req models.GamePlayersRequest) ([]models.LeaderboardRow, error) {
	rows, err := getLeaderboardForPlayers.Query(req.Game, pq.Array(req.Players))
	if err != nil {
		log.Error("Error querying for leaderboard for players in a game: ", "error", err)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Warn("Error closing rows: ", "error", err)
		}
	}(rows)
	var leaderboard = make([]models.LeaderboardRow, 0)
	for rows.Next() {
		var row models.LeaderboardRow
		err = rows.Scan(&row.Game, &row.Player, &row.Position, &row.Score, &row.UnixTimeStamp)
		if err != nil {
			log.Error("Error scanning leaderboard row: ", "error", err)
			return nil, err
		}
		leaderboard = append(leaderboard, row)
	}

	return leaderboard, nil
}
