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

func (d *defaultDatabase) GetSeasons(active bool) ([]models.Season, error) {
	return innerGetSeasons(active)
}

func (d *defaultDatabase) GetChests(season string) ([]models.ChestLocation, error) {
	return innerGetChestLocations(season)
}

func (d *defaultDatabase) GetCurrentChests() ([]models.ChestLocation, error) {
	return innerGetCurrentChestLocations()
}
