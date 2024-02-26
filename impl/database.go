package impl

import (
	"database/sql"
	"fmt"

	"github.com/Fesaa/CubepanionAPI/models"
	_ "github.com/lib/pq"
)

type databaseProviderImpl struct {
	backingDB *sql.DB
}

func newDatabaseProvider(connURL string) (models.DatabaseProvider, error) {
	backingDB, err := sql.Open("postgres", connURL)
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %v", err)
	}

	err = loadQueries(backingDB)
	if err != nil {
		return nil, fmt.Errorf("Error loading queries: %v", err)
	}

	return &databaseProviderImpl{backingDB: backingDB}, nil
}

func (db *databaseProviderImpl) GetBackingDB() *sql.DB {
	return db.backingDB
}

func (db *databaseProviderImpl) GetSeasons(active bool) ([]models.Season, error) {
	return innerGetSeasons(active)
}

func (db *databaseProviderImpl) GetCurrentChestLocations() ([]models.ChestLocation, error) {
	return innerGetCurrentChestLocations()
}

func (db *databaseProviderImpl) GetChestLocations(season string) ([]models.ChestLocation, error) {
	return innerGetChestLocations(season)
}

func (db *databaseProviderImpl) GetEggWarsMaps() ([]models.EggWarsMap, error) {
	return innerGetEggWarsMaps()
}

func (db *databaseProviderImpl) GetEggWarsMap(uniqueName string) (*models.EggWarsMap, error) {
	return innerGetEggWarsMap(uniqueName)
}

func (db *databaseProviderImpl) GetGames(active bool) ([]models.Game, error) {
	return innerGetGames(active)
}

func (db *databaseProviderImpl) GetLeaderboard(game string) ([]models.LeaderboardRow, error) {
	return db.GetLeaderboardBounded(game, 0, 200)
}

func (db *databaseProviderImpl) GetLeaderboardBounded(game string, start int, end int) ([]models.LeaderboardRow, error) {
	return innerGetLeaderboardBounded(game, start, end)
}

func (db *databaseProviderImpl) GetLeaderboardForPlayer(player string) ([]models.LeaderboardRow, error) {
	return innerGetLeaderboardForPlayer(db.backingDB, player)
}
