package models

type LeaderboardSubmission struct {
	Uuid          string           `json:"uuid" validate:"nonzero,regexp=^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"`
	UnixTimeStamp uint64           `json:"unix_time_stamp" validate:"nonzero"`
	Game          string           `json:"game" validate:"nonzero,regexp=[a-zA-Z0-9_ ]"`
	Entries       []LeaderboardRow `json:"entries" validate:"nonzero,len=200"`
}

type GamePlayersRequest struct {
	Game    string   `json:"game" validate:"nonzero,regexp=[a-zA-Z0-9_ ]"`
	Players []string `json:"players" validate:"nonzero"`
}
