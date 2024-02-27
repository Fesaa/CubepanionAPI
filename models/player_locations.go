package models

type PlayerLocationProvider interface {
	// GetPlayerLocation returns the location of the player with the given UUID.
	GetPlayerLocation(uuid string) (*Location, error)

	// SetPlayerLocation sets the location of the player with the given UUID.
	SetPlayerLocation(uuid string, location Location) error

	// RemovePlayerLocation removes the location of the player with the given UUID.
	RemovePlayerLocation(uuid string) error

	// GetSharedPlayers returns the UUIDs of all players who have shared their location with the player with the given UUID.
	GetSharedPlayers(uuid string) ([]string, error)
}

type Location struct {
	Current    string `json:"current"`
	Previous   string `json:"previous"`
	InPreLobby bool   `json:"inPreLobby"`
}
