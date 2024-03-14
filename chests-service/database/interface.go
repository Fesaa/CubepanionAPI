package database

import "github.com/Fesaa/CubepanionAPI/core/models"

type Database interface {
	GetSeasons(active bool) ([]models.Season, error)
	GetChests(season string) ([]models.ChestLocation, error)
	GetCurrentChests() ([]models.ChestLocation, error)
}
