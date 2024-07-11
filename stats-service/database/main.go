package database

import (
	"database/sql"
	"errors"
	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/core/log"
	"github.com/Fesaa/CubepanionAPI/core/models"
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

func (d *defaultDatabase) GetStatsForGame(game string) (*models.GameStat, error) {
	var gs models.GameStat
	err := getGameStat.QueryRow(game).Scan(&gs.Game, &gs.PlayerCount, &gs.UnixTimeStamp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &gs, nil
}

func (d *defaultDatabase) GetAllStats() ([]models.GameStat, error) {
	rows, err := getAllStats.Query()
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Warn("Error closing rows: ", "error", err)
		}
	}(rows)

	stats := make([]models.GameStat, 0)
	for rows.Next() {
		var gs models.GameStat
		if err = rows.Scan(&gs.Game, &gs.PlayerCount, &gs.UnixTimeStamp); err != nil {
			return nil, err
		}
		stats = append(stats, gs)
	}

	return stats, nil
}
