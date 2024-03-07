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

func GetPlayerLocation(uuid string) (*models.Location, error) {
	row := getPlayerLocation.QueryRow(uuid)
	var location models.Location
	err := row.Scan(&location.Current, &location.Previous, &location.InPreLobby)
	if err != nil {
		return nil, err
	}

	return &location, nil
}

func SetPlayerLocation(uuid string, location models.Location) error {
	_, err := setPlayerLocation.Exec(uuid, location.Current, location.Previous, location.InPreLobby)
	return err
}

func RemovePlayerLocation(uuid string) error {
	_, err := removePlayerLocation.Exec(uuid)
	return err
}

func GetSharedPlayers(uuid string) ([]string, error) {
	rows, err := getSharedPlayers.Query(uuid)
	if err != nil {
		return nil, err
	}

	var sharedPlayers []string
	for rows.Next() {
		var sharedPlayer string
		err = rows.Scan(&sharedPlayer)
		if err != nil {
			return nil, err
		}
		sharedPlayers = append(sharedPlayers, sharedPlayer)
	}

	return sharedPlayers, nil
}
