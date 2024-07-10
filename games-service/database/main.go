package database

import (
	"database/sql"
	"errors"
	"github.com/Fesaa/CubepanionAPI/core/log"
	"strings"

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

func (d *defaultDatabase) GetGames(active bool) ([]models.Game, error) {
	var rows *sql.Rows
	var err error

	if active {
		rows, err = getActiveGames.Query()
	} else {
		rows, err = getGames.Query()
	}

	if err != nil {
		log.Error("Error querying for games: ", "error", err)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Warn("Error closing rows: ", "error", err)
		}
	}(rows)
	var games []models.Game = make([]models.Game, 0)
	for rows.Next() {
		var game models.Game
		var aliases string
		err = rows.Scan(&game.Game, &game.Active, &game.DisplayName, &aliases, &game.ScoreType)
		if err != nil {
			log.Error("Error scanning game: ", "error", err)
			return nil, err
		}
		if aliases == "" {
			game.Aliases = []string{}
		} else {
			game.Aliases = strings.Split(aliases, ",")
		}
		games = append(games, game)
	}

	return games, nil
}

func (d *defaultDatabase) GetGame(s string) (string, error) {
	row := getGame.QueryRow(s)

	var game string
	err := row.Scan(&game)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		log.Error("Error scanning game: ", "error", err)
		return "", err
	}

	return game, nil
}
