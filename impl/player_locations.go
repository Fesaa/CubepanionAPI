package impl

import (
	"database/sql"

	"github.com/Fesaa/CubepanionAPI/models"
)

type playerLocationImpl struct {
	db *sql.DB
}

func newPlayerLocation(db *sql.DB) (models.PlayerLocationProvider, error) {
	err := loadPlayerLocationQueries(db)
	if err != nil {
		return nil, err
	}

	return &playerLocationImpl{db: db}, nil
}

func (pl *playerLocationImpl) GetPlayerLocation(uuid string) (*models.Location, error) {
	row := getPlayerLocation.QueryRow(uuid)
	var location models.Location
	err := row.Scan(&location.Current, &location.Previous, &location.InPreLobby)
	if err != nil {
		return nil, err
	}

	return &location, nil
}

func (pl *playerLocationImpl) SetPlayerLocation(uuid string, location models.Location) error {
	_, err := setPlayerLocation.Exec(uuid, location.Current, location.Previous, location.InPreLobby)
	return err
}

func (pl *playerLocationImpl) RemovePlayerLocation(uuid string) error {
	_, err := removePlayerLocation.Exec(uuid)
	return err
}

func (pl *playerLocationImpl) GetSharedPlayers(uuid string) ([]string, error) {
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
