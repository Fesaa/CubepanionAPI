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

func (d *defaultDatabase) GetPlayerLocation(uuid string) (*models.Location, error) {
	row := getPlayerLocation.QueryRow(uuid)
	var location models.Location
	err := row.Scan(&location.Current, &location.Previous, &location.InPreLobby)
	if err != nil {
		return nil, err
	}

	return &location, nil
}

func (d *defaultDatabase) SetPlayerLocation(uuid string, location models.Location) error {
	_, err := setPlayerLocation.Exec(uuid, location.Current, location.Previous, location.InPreLobby)
	return err
}

func (d *defaultDatabase) RemovePlayerLocation(uuid string) error {
	_, err := removePlayerLocation.Exec(uuid)
	return err
}

func (d *defaultDatabase) GetSharedPlayers(uuid string) ([]string, error) {
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

func (d *defaultDatabase) SetProtocolVersion(uuid string, version int) error {
	_, err := setProtocolVersion.Exec(uuid, version)
	return err
}

func (d *defaultDatabase) SetGameStat(stats models.GameStat, uuid string) error {
	_, err := setGameStat.Exec(stats.UnixTimeStamp, stats.Game, stats.PlayerCount, uuid)
	return err
}
