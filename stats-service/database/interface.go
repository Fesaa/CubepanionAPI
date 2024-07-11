package database

import "github.com/Fesaa/CubepanionAPI/core/models"

type Database interface {
	GetStatsForGame(string) (*models.GameStat, error)
	GetAllStats() ([]models.GameStat, error)
}
