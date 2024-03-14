package database

import "github.com/Fesaa/CubepanionAPI/core/models"

type Database interface {
	GetLeaderboard(game string) ([]models.LeaderboardRow, error)
	GetLeaderboardBounded(game string, start, end int) ([]models.LeaderboardRow, error)
	GetLeaderboardForPlayer(player string) ([]models.LeaderboardRow, error)
	SubmitLeaderboard(req models.LeaderboardSubmission) error
}
