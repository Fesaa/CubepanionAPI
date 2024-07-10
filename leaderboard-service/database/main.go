package database

import (
	"database/sql"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/core/models"
	_ "github.com/lib/pq"
)

type defaultDatabase struct {
	db *sql.DB
}

func Connect(d core.DatabaseConfig) (Database, error) {
	db, err := sql.Open("postgres", d.AsConnectionString())
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

	return &defaultDatabase{db: db}, nil
}

func (d *defaultDatabase) GetLeaderboard(game string) ([]models.LeaderboardRow, error) {
	return innerGetLeaderboardBounded(game, 0, models.LEADERBOARD_SIZE)
}

func (d *defaultDatabase) GetLeaderboardBounded(game string, start, end int) ([]models.LeaderboardRow, error) {
	return innerGetLeaderboardBounded(game, start, end)
}

func (d *defaultDatabase) GetLeaderboardForPlayer(player string) ([]models.LeaderboardRow, error) {
	return innerGetLeaderboardForPlayer(player)
}

func (d *defaultDatabase) SubmitLeaderboard(req models.LeaderboardSubmission) error {
	return innerInsertLeaderboards(d.db, req)
}

func (d *defaultDatabase) GetLeaderboardForPlayers(req models.GamePlayersRequest) ([]models.LeaderboardRow, error) {
	return innerGetLeaderboardForPlayers(req)
}
