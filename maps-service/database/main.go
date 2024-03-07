package database

import (
	"database/sql"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/core/models"
	_ "github.com/lib/pq"
)

func Connect(d core.DatabaseConfig) (*sql.DB, error) {
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

	return db, nil
}

func GetMap(name string) (*models.EggWarsMap, error) {
	return innerGetEggWarsMap(name)
}

func GetAllMaps() ([]models.EggWarsMap, error) {
	return innerGetEggWarsMaps()
}
