package impl

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Fesaa/CubepanionAPI/models"
)

type gameSubmission struct {
	game           string
	lastSubmission int64
}

func (g gameSubmission) Value() (driver.Value, error) {
	return fmt.Sprintf("('%s', %d)", g.game, g.lastSubmission), nil
}

func generatePlayerSQL(player string, lsu []gameSubmission) string {
	s := fmt.Sprintf("SELECT * FROM leaderboards WHERE UPPER(player) = UPPER('%s')", player)
	if len(lsu) > 0 {
		subMap := ""
		for _, sub := range lsu {
			subMap += fmt.Sprintf("('%s', %d), ", sub.game, sub.lastSubmission)
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
	gameDisplayName, err := getGameDisplayName(game)
	if err != nil {
		slog.Error("Error getting game display name: %v", err)
		return nil, err
	}

	unixRow := getLastSubmission.QueryRow(gameDisplayName)
	var lastSubmission int64
	err = unixRow.Scan(&lastSubmission)
	if err != nil {
		slog.Error("Error scanning last submission: %v", err)
		return nil, err
	}

	rows, err := getLeaderboard.Query(lastSubmission, gameDisplayName, start, end)
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

func innerGetGames(active bool) ([]models.Game, error) {
	var rows *sql.Rows
	var err error

	if active {
		rows, err = getActiveGames.Query()
	} else {
		rows, err = getGames.Query()
	}

	if err != nil {
		slog.Error("Error querying for games: %v", err)
		return nil, err
	}

	defer rows.Close()
	var games []models.Game = make([]models.Game, 0)
	for rows.Next() {
		var game models.Game
		var aliases string
		err = rows.Scan(&game.Game, &game.Active, &game.DisplayName, &aliases, &game.ScoreType)
		if err != nil {
			slog.Error("Error scanning game: %v", err)
			return nil, err
		}
		game.Aliases = strings.Split(aliases, ",")
		games = append(games, game)
	}

	return games, nil
}

func getGameDisplayName(game string) (string, error) {
	gameRow := getGame.QueryRow(game)
	var gameDisplayName string
	err := gameRow.Scan(&gameDisplayName)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("No game with name %s", game)
		}
		return "", err
	}

	return gameDisplayName, nil
}
