package database

import (
	"database/sql"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/core/models"
	_ "github.com/lib/pq"
)

var db *sql.DB

func Connect(d core.DatabaseConfig) (*sql.DB, error) {
	var err error
	db, err = sql.Open("postgres", d.AsConnectionString())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = load(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetLeaderboard(game string) ([]models.LeaderboardRow, error) {
	return innerGetLeaderboardBounded(game, 0, models.LEADERBOARD_SIZE)
}

func GetLeaderboardBounded(game string, start, end int) ([]models.LeaderboardRow, error) {
	return innerGetLeaderboardBounded(game, start, end)
}

func GetLeaderboardForPlayer(player string) ([]models.LeaderboardRow, error) {
	return innerGetLeaderboardForPlayer(db, player)
}

func SubmitLeaderboard(req models.LeaderboardSubmission) error {
	return innerInsertLeaderboards(db, req)
}
