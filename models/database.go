package models

import "database/sql"

type DatabaseProvider interface {
	GetBackingDB() *sql.DB

	// GetSeasons returns all seasons if active is true, otherwise it returns all seasons
	GetSeasons(active bool) ([]Season, error)
	// GetCurrentChestLocations returns all chest locations for the current season
	GetCurrentChestLocations() ([]ChestLocation, error)
	// GetChestLocations returns all chest locations for the given season
	GetChestLocations(season string) ([]ChestLocation, error)

	// GetEggWarsMaps returns all egg wars maps
	GetEggWarsMaps() ([]EggWarsMap, error)
	// GetEggWarsMap returns the egg wars map with the given unique name
	GetEggWarsMap(uniqueName string) (*EggWarsMap, error)

	// GetGames returns all games if active is true, otherwise it returns all games
	GetGames(active bool) ([]Game, error)

	// GetLeaderboard returns the top 200 players for the given game. This may be the UniqueName, DisplayName or an Alias
	GetLeaderboard(game string) ([]LeaderboardRow, error)
	// GetLeaderboardBounded returns the top players for the given game between the start and end
	GetLeaderboardBounded(game string, start int, end int) ([]LeaderboardRow, error)
	// GetLeaderboardForPlayer returns the leaderboard for the given player
	GetLeaderboardForPlayer(player string) ([]LeaderboardRow, error)
}
