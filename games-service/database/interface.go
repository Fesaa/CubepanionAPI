package database

import "github.com/Fesaa/CubepanionAPI/core/models"

type Database interface {
	GetGames(active bool) ([]models.Game, error)
	GetGame(s string) (string, error)
}
