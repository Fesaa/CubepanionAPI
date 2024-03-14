package database

import "github.com/Fesaa/CubepanionAPI/core/models"

type Database interface {
	GetMap(name string) (*models.EggWarsMap, error)
	GetAllMaps() ([]models.EggWarsMap, error)
}
