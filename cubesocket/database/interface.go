package database

import "github.com/Fesaa/CubepanionAPI/core/models"

type Database interface {
	GetPlayerLocation(uuid string) (*models.Location, error)
	SetPlayerLocation(uuid string, location models.Location) error
	RemovePlayerLocation(uuid string) error
	GetSharedPlayers(uuid string) ([]string, error)
	SetProtocolVersion(uuid string, version int) error
	SetGameStat(stats models.GameStat, uuid string) error
}
