package database

import "github.com/Fesaa/CubepanionAPI/core/models"

type Database interface {
	GetAllMaps() ([]models.GameMap, error)
}
